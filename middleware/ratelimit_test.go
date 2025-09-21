package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rushairer/gouno/middleware"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter(t *testing.T) {
	// 创建一个限制为每分钟3次的限频器用于测试
	limiter := middleware.NewRateLimiter(3, time.Minute)

	// 测试正常请求
	assert.True(t, limiter.IsAllowed("192.168.1.1"))
	assert.True(t, limiter.IsAllowed("192.168.1.1"))
	assert.True(t, limiter.IsAllowed("192.168.1.1"))

	// 第4次请求应该被拒绝
	assert.False(t, limiter.IsAllowed("192.168.1.1"))

	// 不同IP应该有独立的限制
	assert.True(t, limiter.IsAllowed("192.168.1.2"))
}

func TestRateLimiterRemainingRequests(t *testing.T) {
	limiter := middleware.NewRateLimiter(5, time.Minute)
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

	// 创建测试路由
	r := gin.New()
	r.Use(middleware.RateLimitMiddleware(2, time.Minute)) // 每分钟2次
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// 第1次请求 - 应该成功
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "60", w1.Header().Get("X-RateLimit-Limit"))

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
	assert.Equal(t, "60", w3.Header().Get("Retry-After"))
}

func TestRateLimitMiddlewareDifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.RateLimitMiddleware(1, time.Minute)) // 每分钟1次
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
	// 使用很短的时间窗口进行测试
	limiter := middleware.NewRateLimiter(2, 100*time.Millisecond)
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

// 基准测试
func BenchmarkRateLimiter(b *testing.B) {
	limiter := middleware.NewRateLimiter(1000, time.Minute)
	ip := "192.168.1.1"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.IsAllowed(ip)
	}
}

func BenchmarkRateLimiterConcurrent(b *testing.B) {
	limiter := middleware.NewRateLimiter(1000, time.Minute)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		ip := "192.168.1.1"
		for pb.Next() {
			limiter.IsAllowed(ip)
		}
	})
}
