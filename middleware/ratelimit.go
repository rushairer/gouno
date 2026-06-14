package middleware

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rushairer/gouno"
)

const defaultMaxVisitors = 10000

// RateLimiter IP 限频器
type RateLimiter struct {
	mu          sync.RWMutex
	visitors    map[string]*Visitor
	limit       int           // 每分钟允许的请求数
	window      time.Duration // 时间窗口
	maxVisitors int           // visitors map 最大条目数，防止内存膨胀
}

// Visitor 访问者信息
type Visitor struct {
	requests []time.Time // 请求时间戳列表
	mu       sync.Mutex
}

// NewRateLimiter 创建新的限频器，ctx 取消时停止后台清理 goroutine
func NewRateLimiter(ctx context.Context, limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors:    make(map[string]*Visitor),
		limit:       limit,
		window:      window,
		maxVisitors: defaultMaxVisitors,
	}

	go rl.cleanupVisitors(ctx)

	return rl
}

// SetMaxVisitors 设置 visitors map 的最大条目数。超过上限时新 IP 将被拒绝直到清理完成。
// 传入 0 则使用默认值 (10000)。
func (rl *RateLimiter) SetMaxVisitors(max int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if max <= 0 {
		rl.maxVisitors = defaultMaxVisitors
	} else {
		rl.maxVisitors = max
	}
}

// IsAllowed 检查是否允许请求
func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mu.Lock()
	visitor, exists := rl.visitors[ip]
	if !exists {
		// 超过最大条目数时，触发激进清理并拒绝新 IP
		if len(rl.visitors) >= rl.maxVisitors {
			rl.evictIdleVisitors()
			if len(rl.visitors) >= rl.maxVisitors {
				rl.mu.Unlock()
				return false
			}
		}
		visitor = &Visitor{
			requests: make([]time.Time, 0),
		}
		rl.visitors[ip] = visitor
	}
	rl.mu.Unlock()

	visitor.mu.Lock()
	defer visitor.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// 移除过期的请求记录
	validRequests := make([]time.Time, 0)
	for _, reqTime := range visitor.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	visitor.requests = validRequests

	// 检查是否超过限制
	if len(visitor.requests) >= rl.limit {
		return false
	}

	// 添加当前请求
	visitor.requests = append(visitor.requests, now)
	return true
}

// cleanupVisitors 定期清理过期的访问者，ctx 取消时退出
func (rl *RateLimiter) cleanupVisitors(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			cutoff := now.Add(-rl.window * 2) // 保留2倍窗口时间的数据

			for ip, visitor := range rl.visitors {
				visitor.mu.Lock()
				// 如果访问者在2倍窗口时间内没有请求，则删除
				if len(visitor.requests) == 0 || visitor.requests[len(visitor.requests)-1].Before(cutoff) {
					delete(rl.visitors, ip)
				}
				visitor.mu.Unlock()
			}
			rl.mu.Unlock()
		}
	}
}

// evictIdleVisitors 移除所有当前无活跃请求的访问者以释放 map 空间。
// 包括 requests 为空的访问者，以及所有请求均已过期的访问者。
// 调用方必须持有 rl.mu 写锁。
func (rl *RateLimiter) evictIdleVisitors() {
	now := time.Now()
	cutoff := now.Add(-rl.window)
	for ip, visitor := range rl.visitors {
		// 已持有 rl.mu 写锁，无其他 goroutine 可访问此 visitor 指针，
		// 因此无需加 visitor.mu。
		if len(visitor.requests) == 0 {
			delete(rl.visitors, ip)
			continue
		}
		// 如果最后一个（最新的）请求也已过期，说明该访问者不再活跃
		if visitor.requests[len(visitor.requests)-1].Before(cutoff) {
			visitor.requests = visitor.requests[:0] // 释放底层数组引用
			delete(rl.visitors, ip)
		}
	}
}

// GetRemainingRequests 获取剩余请求数
func (rl *RateLimiter) GetRemainingRequests(ip string) int {
	rl.mu.RLock()
	visitor, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists {
		return rl.limit
	}

	visitor.mu.Lock()
	defer visitor.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// 计算有效请求数
	validCount := 0
	for _, reqTime := range visitor.requests {
		if reqTime.After(cutoff) {
			validCount++
		}
	}

	remaining := rl.limit - validCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetResetTime 获取限制重置时间
func (rl *RateLimiter) GetResetTime(ip string) time.Time {
	rl.mu.RLock()
	visitor, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists || len(visitor.requests) == 0 {
		return time.Now()
	}

	visitor.mu.Lock()
	defer visitor.mu.Unlock()

	// 找到最早的有效请求
	now := time.Now()
	cutoff := now.Add(-rl.window)

	for _, reqTime := range visitor.requests {
		if reqTime.After(cutoff) {
			return reqTime.Add(rl.window)
		}
	}

	return time.Now()
}

// RateLimitMiddleware 创建限频中间件
func RateLimitMiddleware(ctx context.Context, limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(ctx, limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 检查是否允许请求
		if !limiter.IsAllowed(ip) {
			resetTime := limiter.GetResetTime(ip)

			// 设置响应头
			c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))
			c.Header("Retry-After", strconv.Itoa(int(window.Seconds())))

			c.JSON(http.StatusTooManyRequests, gouno.NewErrorResponse(http.StatusTooManyRequests, "too many requests"))
			c.Abort()
			return
		}

		// 设置成功响应的限频头
		remaining := limiter.GetRemainingRequests(ip)
		resetTime := limiter.GetResetTime(ip)

		// 确保所有头部值都是有效的字符串，避免空字节
		limitStr := strconv.Itoa(limit)
		remainingStr := strconv.Itoa(remaining)
		resetStr := resetTime.Format(time.RFC3339)

		c.Header("X-RateLimit-Limit", limitStr)
		c.Header("X-RateLimit-Remaining", remainingStr)
		c.Header("X-RateLimit-Reset", resetStr)

		c.Next()
	}
}

// IPRateLimitMiddleware 默认的 IP 限频中间件 (每分钟60次)
func IPRateLimitMiddleware(ctx context.Context) gin.HandlerFunc {
	return RateLimitMiddleware(ctx, 60, time.Minute)
}
