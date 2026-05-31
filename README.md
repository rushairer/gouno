# gouno

[English](#english) | [中文](#中文)

---

## English

### What is gouno?

gouno is a **lightweight Go web project launcher**. It scaffolds the project structure, startup flow, web layer, and response format — the boilerplate you'd rewrite for every microservice — so you can focus on business logic from the first line.

gouno is **not a framework**. It doesn't bundle a database driver, cache client, or message queue. Your tech stack is your choice. gouno handles the rest.

```
What gouno does                    What gouno doesn't do
├── Project structure (DDD)        ├── Database (pgx? gorm? ent?)
├── CLI + config (Cobra + Viper)   ├── Cache (redis? memcached?)
├── Web engine (Gin + middleware)   ├── Messaging (kafka? rabbitmq?)
├── Response format (unified JSON)  └── Auth (JWT? OAuth? session?)
└── Code generator + templates
```

### Quick Start

```bash
# Install gouno-cli
go install github.com/rushairer/gouno-cli@latest

# Create a new project
gouno-cli new my-service -m github.com/you/my-service

# Build and run
cd my-service
go mod tidy
make dev
# → http://localhost:8080
```

Your service is running. From here, write business code.

### Code Generation

Generate a complete DDD module (domain + repository + service) in one command:

```bash
gouno gen suite user
```

```
internal/
├── domain/user.go         ← Entity
├── repository/user.go     ← Data access interface
└── service/user.go        ← Business logic interface
```

Generate individual pieces:

```bash
gouno gen domain order
gouno gen repository order
gouno gen service order
gouno gen controller order
gouno gen task send_email
```

### Template Sets

Different teams, different code styles. Template sets let you customize what `gouno gen` produces — without touching gouno's source code.

```bash
# Install a community or company template set
gouno-cli template install gorm https://github.com/myorg/gouno-template-gorm

# Use it when creating a project
gouno-cli new order-service --template-set gorm -m github.com/myorg/order-service

# Now gen suite uses your gorm templates
cd order-service
gouno gen suite user   # → generates gorm-style domain/repository/service
```

Manage template sets:

```bash
gouno-cli template list              # List installed sets
gouno-cli template install <n> <url> # Install from git or local path
gouno-cli template remove <n>        # Remove a set
```

**How it works:**

```
Search priority:
1. --template-set flag (highest)
2. .gouno.yaml in project root
3. Built-in default templates

Template locations:
~/.gouno/templates/
├── default/      ← Ships with gouno-template
├── gorm/         ← Your custom set
└── my-company/   ← Your company set
```

Create your own template set — just a directory with `.tmpl` files:

```
my-template-set/
├── domain.tmpl
├── repository.tmpl
├── service.tmpl
├── controller.tmpl
└── task.tmpl
```

Each `.tmpl` file uses `%s` as the struct name placeholder:

```
package domain

import "time"

type %s struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"size:255"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Project Structure

```
my-service/
├── cmd/main.go                      ← Entry point
│   └── gouno/
│       ├── root.go                  ← Cobra root command
│       └── web.go                   ← Web server (Gin + middleware + graceful shutdown)
├── config/                          ← Multi-environment YAML config
│   ├── development.yaml
│   ├── test.yaml
│   ├── production.yaml
│   ├── config.go                    ← Config struct definitions
│   └── config_manager.go            ← Thread-safe config loader
├── controller/                      ← HTTP handlers (generated)
├── internal/                        ← Business modules (DDD)
│   ├── domain/                      ← Entities
│   ├── repository/                  ← Data access
│   ├── service/                     ← Business logic
│   └── task/                        ← Background tasks
├── middleware/                      ← Custom middleware
├── router/                          ← Route registration
├── .gouno.yaml                      ← Template set config (optional)
├── Makefile                         ← Build/run/dev shortcuts
└── go.mod
```

### API Response Format

All responses follow a unified JSON structure:

```json
// Success (HTTP 200)
{"code": 200, "message": "success", "data": {...}}

// Error (HTTP 4xx/5xx)
{"code": 400, "message": "bad request", "data": null}
{"code": 404, "message": "not found", "data": null}
{"code": 500, "message": "internal server error", "data": null}
```

Usage in handlers:

```go
ctx.JSON(http.StatusOK, gouno.NewSuccessResponse(data))
ctx.JSON(http.StatusBadRequest, gouno.BadRequestResponse)
```

### Middleware

Built-in middleware chain (in order):

| Middleware | Purpose | HTTP Status |
|------------|---------|-------------|
| Logger | Request logging | — |
| Recovery | Panic recovery | 500 |
| Timeout | Request timeout | 408 |
| RateLimit | IP-based rate limiting | 429 |

Add your own middleware in `web.go`:

```go
engine.Use(
    gin.Logger(),
    middleware.RecoveryMiddleware(),
    middleware.TimeoutMiddleware(config.RequestTimeout),
    gounoMiddleware.RateLimitMiddleware(ctx, config.RateLimitPerMinute, time.Minute),
    // your middleware here
)
```

### Philosophy

**gouno is a launcher, not a framework.**

- **Launcher**: Handles startup boilerplate. You own the code. Modify anything.
- **Framework**: Abstracts away details. You follow its rules. Hard to deviate.

gouno gives you a standardized starting point, then gets out of the way. The real customization happens in **template sets** — your team's code style, your tech stack choices, your patterns — all captured as reusable templates.

```
gouno (core)          → Handles: project structure + startup + web layer + response
template sets         → Handles: code style + tech stack + team conventions
your business code    → Handles: actual product logic
```

---

## 中文

### gouno 是什么？

gouno 是一个**轻量级 Go Web 项目启动器**。它为你搭建项目结构、启动流程、Web 层和响应格式——这些每个微服务都要写、但跟业务无关的代码——让你从第一行就专注业务逻辑。

gouno **不是框架**。它不绑定数据库驱动、缓存客户端或消息队列。技术选型由你决定，gouno 只负责其余部分。

### 快速开始

```bash
# 安装 gouno-cli
go install github.com/rushairer/gouno-cli@latest

# 创建新项目
gouno-cli new my-service -m github.com/you/my-service

# 构建运行
cd my-service
go mod tidy
make dev
# → http://localhost:8080
```

### 代码生成

一条命令生成完整的 DDD 模块（domain + repository + service）：

```bash
gouno gen suite user
```

也支持单独生成：

```bash
gouno gen domain order
gouno gen controller auth
gouno gen task send_email
```

### 模板集

不同团队、不同代码风格。模板集让你自定义 `gouno gen` 的输出，无需修改 gouno 源码：

```bash
# 安装模板集
gouno-cli template install gorm https://github.com/myorg/gouno-template-gorm

# 创建项目时指定模板集
gouno-cli new order-service --template-set gorm -m github.com/myorg/order-service

# gen suite 自动使用 gorm 模板
cd order-service
gouno gen suite user
```

创建自定义模板集：只需一个包含 `.tmpl` 文件的目录，`%s` 作为结构体名占位符。

### 设计理念

**gouno 是启动器，不是框架。**

gouno 处理项目启动的标准化工作，然后让开。真正的定制通过**模板集**完成——你的代码风格、技术选型、团队规范——全部封装为可复用的模板。

gouno 负责启动标准化，模板集负责代码风格，你负责业务逻辑。
