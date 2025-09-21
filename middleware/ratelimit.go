package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter IP 限频器
type RateLimiter struct {
	mu       sync.RWMutex
	visitors map[string]*Visitor
	limit    int           // 每分钟允许的请求数
	window   time.Duration // 时间窗口
}

// Visitor 访问者信息
type Visitor struct {
	requests []time.Time // 请求时间戳列表
	mu       sync.Mutex
}

// NewRateLimiter 创建新的限频器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		limit:    limit,
		window:   window,
	}

	// 启动清理 goroutine，定期清理过期的访问者
	go rl.cleanupVisitors()

	return rl
}

// IsAllowed 检查是否允许请求
func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mu.Lock()
	visitor, exists := rl.visitors[ip]
	if !exists {
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

// cleanupVisitors 定期清理过期的访问者
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟清理一次
	defer ticker.Stop()

	for range ticker.C {
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
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 检查是否允许请求
		if !limiter.IsAllowed(ip) {
			resetTime := limiter.GetResetTime(ip)

			// 设置响应头
			c.Header("X-RateLimit-Limit", "60")
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))
			c.Header("Retry-After", "60")

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":    "Rate limit exceeded",
				"message":  "Too many requests. Please try again later.",
				"limit":    limit,
				"window":   window.String(),
				"reset_at": resetTime.Format(time.RFC3339),
			})
			c.Abort()
			return
		}

		// 设置成功响应的限频头
		remaining := limiter.GetRemainingRequests(ip)
		resetTime := limiter.GetResetTime(ip)

		c.Header("X-RateLimit-Limit", "60")
		c.Header("X-RateLimit-Remaining", string(rune(remaining)))
		c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

		c.Next()
	}
}

// IPRateLimitMiddleware 默认的 IP 限频中间件 (每分钟60次)
func IPRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(60, time.Minute)
}
