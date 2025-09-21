# Rate Limit Middleware

IP 限频中间件，用于限制每个 IP 地址的请求频率，防止 API 滥用和 DDoS 攻击。

## 功能特性

- ✅ **IP 级别限频**：基于客户端 IP 地址进行限制
- ✅ **滑动窗口算法**：使用滑动时间窗口，更精确的限频控制
- ✅ **内存存储**：高性能的内存存储，无需外部依赖
- ✅ **自动清理**：定期清理过期的访问记录，防止内存泄漏
- ✅ **标准响应头**：符合 HTTP 标准的限频响应头
- ✅ **并发安全**：支持高并发场景
- ✅ **灵活配置**：可自定义限制次数和时间窗口

## 快速开始

### 基本使用

```go
package main

import (
    "gosso/middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    // 使用默认限频：每分钟 60 次
    r.Use(middleware.IPRateLimitMiddleware())
    
    r.GET("/api/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })
    
    r.Run(":8080")
}
```

### 自定义限频配置

```go
import (
    "time"
    "gosso/middleware"
)

func setupRoutes() {
    r := gin.Default()
    
    // 全局限频：每分钟 100 次
    r.Use(middleware.RateLimitMiddleware(100, time.Minute))
    
    // API 路由组
    api := r.Group("/api/v1")
    {
        // 账户接口：默认限频
        account := api.Group("/account")
        {
            account.POST("/register", handleRegister)
        }
        
        // 登录接口：更严格的限频（每分钟 10 次）
        auth := api.Group("/auth")
        auth.Use(middleware.RateLimitMiddleware(10, time.Minute))
        {
            auth.POST("/login", handleLogin)
            auth.POST("/forgot-password", handleForgotPassword)
        }
        
        // 公开接口：更宽松的限频（每分钟 200 次）
        public := api.Group("/public")
        public.Use(middleware.RateLimitMiddleware(200, time.Minute))
        {
            public.GET("/status", handleStatus)
        }
    }
}
```

## API 参考

### 函数

#### `IPRateLimitMiddleware() gin.HandlerFunc`
创建默认的 IP 限频中间件（每分钟 60 次）。

#### `RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc`
创建自定义配置的限频中间件。

**参数：**
- `limit`：时间窗口内允许的最大请求次数
- `window`：时间窗口大小

**示例：**
```go
// 每分钟 30 次
middleware.RateLimitMiddleware(30, time.Minute)

// 每小时 1000 次
middleware.RateLimitMiddleware(1000, time.Hour)

// 每秒 10 次
middleware.RateLimitMiddleware(10, time.Second)
```

### 响应头

当请求通过限频检查时，会设置以下响应头：

- `X-RateLimit-Limit`：限制的总次数
- `X-RateLimit-Remaining`：剩余可用次数
- `X-RateLimit-Reset`：限制重置时间（RFC3339 格式）

当请求被限频时，会返回 429 状态码并设置：

- `X-RateLimit-Limit`：限制的总次数
- `X-RateLimit-Remaining`：0
- `X-RateLimit-Reset`：限制重置时间
- `Retry-After`：建议重试等待时间（秒）

### 错误响应

当请求被限频时，返回 JSON 格式的错误信息：

```json
{
  "error": "Rate limit exceeded",
  "message": "Too many requests. Please try again later.",
  "limit": 60,
  "window": "1m0s",
  "reset_at": "2023-12-01T10:31:00Z"
}
```

## 配置建议

### 不同场景的推荐配置

| 场景 | 限制次数 | 时间窗口 | 说明 |
|------|----------|----------|------|
| 登录接口 | 5-10 | 1分钟 | 防止暴力破解 |
| 注册接口 | 3-5 | 1分钟 | 防止恶意注册 |
| 发送验证码 | 1-3 | 1分钟 | 防止短信轰炸 |
| 普通 API | 60-100 | 1分钟 | 正常业务使用 |
| 公开接口 | 100-200 | 1分钟 | 查询类接口 |
| 文件上传 | 10-20 | 1分钟 | 资源密集型操作 |

### 性能考虑

- **内存使用**：每个 IP 大约占用 100-200 字节内存
- **清理机制**：每 5 分钟自动清理过期数据
- **并发性能**：支持高并发，使用读写锁优化性能

## 测试

运行测试：

```bash
# 运行所有测试
go test ./middleware

# 运行特定测试
go test ./middleware -run TestRateLimiter

# 运行基准测试
go test ./middleware -bench=.

# 查看测试覆盖率
go test ./middleware -cover
```

## 高级用法

### 自定义 IP 获取逻辑

如果需要自定义 IP 获取逻辑（例如处理代理），可以修改中间件：

```go
func CustomIPRateLimitMiddleware() gin.HandlerFunc {
    limiter := NewRateLimiter(60, time.Minute)
    
    return func(c *gin.Context) {
        // 自定义 IP 获取逻辑
        ip := c.GetHeader("X-Real-IP")
        if ip == "" {
            ip = c.GetHeader("X-Forwarded-For")
        }
        if ip == "" {
            ip = c.ClientIP()
        }
        
        if !limiter.IsAllowed(ip) {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### 基于用户的限频

```go
func UserRateLimitMiddleware() gin.HandlerFunc {
    limiter := NewRateLimiter(100, time.Minute)
    
    return func(c *gin.Context) {
        // 获取用户 ID
        userID := c.GetString("user_id")
        if userID == "" {
            userID = c.ClientIP() // 未登录用户使用 IP
        }
        
        if !limiter.IsAllowed(userID) {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## 注意事项

1. **内存使用**：大量不同 IP 访问会占用较多内存
2. **重启丢失**：服务重启后限频计数会重置
3. **分布式部署**：多实例部署时每个实例独立计数
4. **代理环境**：注意正确获取真实 IP 地址

## 生产环境建议

1. **监控告警**：监控 429 错误率，及时发现异常流量
2. **日志记录**：记录被限频的请求，便于分析
3. **白名单机制**：为可信 IP 设置白名单
4. **Redis 存储**：大规模部署时考虑使用 Redis 存储
5. **动态配置**：支持运行时调整限频参数