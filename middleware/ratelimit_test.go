package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"

	"github.com/rushairer/gouno/middleware"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiter := middleware.NewRateLimiter(ctx, 3, time.Minute)

	assert.True(t, limiter.IsAllowed("192.168.1.1"))
	assert.True(t, limiter.IsAllowed("192.168.1.1"))
	assert.True(t, limiter.IsAllowed("192.168.1.1"))

	// 第4次请求应该被拒绝
	assert.False(t, limiter.IsAllowed("192.168.1.1"))

	// 不同IP应该有独立的限制
	assert.True(t, limiter.IsAllowed("192.168.1.2"))
}

func TestRateLimiterRemainingRequests(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiter := middleware.NewRateLimiter(ctx, 5, time.Minute)
	ip := "192.168.1.1"

	// 初始状态应该有5次剩余
	assert.Equal(t, 5, limiter.GetRemainingRequests(ip))

	// 使用2次后应该剩余3次
	limiter.IsAllowed(ip)
	limiter.IsAllowed(ip)
	assert.Equal(t, 3, limiter.GetRemainingRequests(ip))

	// 用完所有次数后应该剩余0次
	limiter.IsAllowed(ip)
	limiter.IsAllowed(ip)
	limiter.IsAllowed(ip)
	assert.Equal(t, 0, limiter.GetRemainingRequests(ip))
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	limit := 2

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := gin.New()
	r.Use(middleware.RateLimitMiddleware(ctx, limit, time.Minute))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// 第1次请求 - 应该成功
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "2", w1.Header().Get("X-RateLimit-Limit"))

	// 第2次请求 - 应该成功
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	// 第3次请求 - 应该被限制
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "192.168.1.1:12345"
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusTooManyRequests, w3.Code)
	assert.Equal(t, "0", w3.Header().Get("X-RateLimit-Remaining"))
}

func TestRateLimitMiddlewareDifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := gin.New()
	r.Use(middleware.RateLimitMiddleware(ctx, 1, time.Minute))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// IP1 的第1次请求 - 应该成功
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// IP1 的第2次请求 - 应该被限制
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)

	// IP2 的第1次请求 - 应该成功（不同IP独立计算）
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "192.168.1.2:12345"
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}

func TestRateLimiterTimeWindow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 使用很短的时间窗口进行测试
	limiter := middleware.NewRateLimiter(ctx, 2, 100*time.Millisecond)
	ip := "192.168.1.1"

	// 使用完所有请求
	assert.True(t, limiter.IsAllowed(ip))
	assert.True(t, limiter.IsAllowed(ip))
	assert.False(t, limiter.IsAllowed(ip))

	// 等待时间窗口过期
	time.Sleep(150 * time.Millisecond)

	// 现在应该可以再次请求
	assert.True(t, limiter.IsAllowed(ip))
}

func TestRateLimiterContextCancel(t *testing.T) {
	before := runtime.NumGoroutine()

	ctx, cancel := context.WithCancel(context.Background())
	_ = middleware.NewRateLimiter(ctx, 10, time.Minute)

	// 等待 goroutine 启动
	time.Sleep(10 * time.Millisecond)
	afterCreate := runtime.NumGoroutine()
	assert.GreaterOrEqual(t, afterCreate, before, "goroutine should be created")

	// 取消 context 后 goroutine 应该退出
	cancel()
	time.Sleep(10 * time.Millisecond)
	afterCancel := runtime.NumGoroutine()
	assert.LessOrEqual(t, afterCancel, afterCreate, "goroutine should exit after context cancel")
}

func TestRateLimiterContextCancelMultiple(t *testing.T) {
	before := runtime.NumGoroutine()

	// 创建多个限频器，取消后都应该释放 goroutine
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		_ = middleware.NewRateLimiter(ctx, 10, time.Minute)
		cancel()
	}

	time.Sleep(50 * time.Millisecond)
	after := runtime.NumGoroutine()
	// 允许少量 goroutine 波动，但不应累积 5 个以上
	assert.LessOrEqual(t, after-before+2, 3, "cancelled limiters should not leak goroutines")
}

// 基准测试
func BenchmarkRateLimiter(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiter := middleware.NewRateLimiter(ctx, 1000, time.Minute)
	ip := "192.168.1.1"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.IsAllowed(ip)
	}
}

func BenchmarkRateLimiterConcurrent(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiter := middleware.NewRateLimiter(ctx, 1000, time.Minute)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		ip := "192.168.1.1"
		for pb.Next() {
			limiter.IsAllowed(ip)
		}
	})
}
