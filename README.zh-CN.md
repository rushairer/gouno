# gouno

[English](./README.md)

---

**轻量级 Go Web 项目启动器**。它为你搭建项目结构、启动流程、Web 层和响应格式——这些每个微服务都要写、但跟业务无关的代码——让你从第一行就专注业务逻辑。

gouno **不是框架**。它不绑定数据库驱动、缓存客户端或消息队列。技术选型由你决定，gouno 只负责其余部分。

```
gouno 负责的                         gouno 不负责的
├── 项目结构（DDD）                   ├── 数据库（pgx? gorm? ent?）
├── CLI + 配置（Cobra + Viper）       ├── 缓存（redis? memcached?）
├── Web 引擎（Gin + 中间件）          ├── 消息队列（kafka? rabbitmq?）
├── 响应格式（统一 JSON）             └── 认证（JWT? OAuth? session?）
└── 代码生成器 + 模板集
```

## 快速开始

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

服务已启动，开始写业务代码。

## 代码生成

一条命令生成完整的 DDD 模块（domain + repository + service）：

```bash
gouno gen suite user
```

```
internal/
├── domain/user.go         ← 实体
├── repository/user.go     ← 数据访问接口
└── service/user.go        ← 业务逻辑接口
```

单独生成：

```bash
gouno gen domain order
gouno gen repository order
gouno gen service order
gouno gen controller order
gouno gen task send_email
```

快捷别名：

```bash
gouno gen d order   # domain
gouno gen r order   # repository
gouno gen s order   # service
gouno gen c order   # controller
gouno gen t order   # task
```

选项：

```bash
gouno gen suite user --force          # 覆盖已有文件
gouno gen suite user --path ./custom  # 自定义输出目录
```

## 模板集

不同团队、不同代码风格。模板集让你自定义 `gouno gen` 的输出，无需修改 gouno 源码。

### 安装和使用

```bash
# 安装社区或公司模板集
gouno-cli template install gorm https://github.com/myorg/gouno-template-gorm

# 创建项目时指定模板集
gouno-cli new order-service --template-set gorm -m github.com/myorg/order-service

# gen suite 自动使用 gorm 模板
cd order-service
gouno gen suite user   # → 生成 gorm 风格的 domain/repository/service
```

### 管理模板集

```bash
gouno-cli template list                        # 列出已安装的模板集
gouno-cli template install <name> <git-url>    # 从 git 安装
gouno-cli template install <name> <local-path> # 从本地目录安装
gouno-cli template remove <name>               # 删除模板集
gouno-cli template install <name> <url> --force # 覆盖已有模板集
```

### 模板搜索机制

```
搜索优先级：
1. --template-set 命令行参数（最高）
2. 项目根目录 .gouno.yaml 配置
3. 内置默认模板（兜底）

安装位置：
~/.gouno/templates/
├── default/      ← gouno-template 自带
├── gorm/         ← 自定义模板集
└── my-company/   ← 公司模板集
```

### 创建自定义模板集

模板集就是一个包含 `.tmpl` 文件的目录：

```
my-template-set/
├── domain.tmpl
├── repository.tmpl
├── service.tmpl
├── controller.tmpl
└── task.tmpl
```

每个 `.tmpl` 文件用 `%s` 作为结构体名占位符（5 处）：

```go
package domain

import "time"

type %s struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"size:255"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func New%s() *%s {
    return &%s{}
}
```

把模板集发布为 git 仓库，其他人一条命令即可安装：

```bash
gouno-cli template install my-set https://github.com/myorg/my-template-set
```

## 项目结构

```
my-service/
├── cmd/
│   ├── main.go              ← 入口
│   └── gouno/
│       ├── root.go          ← Cobra 根命令
│       └── web.go           ← Web 服务启动
├── config/
│   ├── development.yaml     ← 开发环境配置
│   ├── test.yaml            ← 测试环境配置
│   ├── production.yaml      ← 生产环境配置（敏感信息用环境变量）
│   ├── config.go            ← 配置结构体定义
│   └── config_manager.go    ← 线程安全的配置加载器
├── controller/              ← HTTP 处理器
├── internal/
│   ├── domain/              ← 实体
│   ├── repository/          ← 数据访问接口
│   ├── service/             ← 业务逻辑
│   └── task/                ← 后台任务
├── middleware/              ← 自定义中间件
├── router/                  ← 路由注册
├── utility/                 ← 公共工具
├── .gouno.yaml              ← 模板集配置（可选）
├── Makefile                 ← 构建/运行/开发快捷命令
├── .air.toml                ← 热重载配置
└── go.mod
```

## API 响应格式

所有响应使用统一的 JSON 结构，`code` 字段与 HTTP 状态码一致：

```json
// 成功（HTTP 200）
{"code": 200, "message": "success", "data": {...}}

// 错误（HTTP 400/401/403/404/408/429/500）
{"code": 400, "message": "bad request", "data": null}
{"code": 404, "message": "not found", "data": null}
{"code": 500, "message": "internal server error", "data": null}
```

在处理器中使用：

```go
// 成功
ctx.JSON(http.StatusOK, gouno.NewSuccessResponse(data))

// 错误（预定义）
ctx.JSON(http.StatusBadRequest, gouno.BadRequestResponse)
ctx.JSON(http.StatusNotFound, gouno.NotFoundResponse)
ctx.JSON(http.StatusInternalServerError, gouno.InternalServerErrorResponse)

// 错误（自定义消息）
ctx.JSON(http.StatusBadRequest, gouno.NewErrorResponse(http.StatusBadRequest, "邮箱已存在"))
```

## 中间件

内置中间件链（按顺序执行）：

| 中间件 | 用途 | 错误时的 HTTP 状态码 |
|--------|------|---------------------|
| Logger | 请求/响应日志 | — |
| Recovery | panic 恢复 | 500 |
| Timeout | 请求超时 | 408 |
| RateLimit | 基于 IP 的滑动窗口限流 | 429 |

在 `web.go` 中添加自定义中间件：

```go
engine.Use(
    gin.Logger(),
    middleware.RecoveryMiddleware(),
    middleware.TimeoutMiddleware(config.RequestTimeout),
    gounoMiddleware.RateLimitMiddleware(ctx, config.RateLimitPerMinute, time.Minute),
    // 你的中间件
)
```

## 配置管理

基于 Viper 的多环境 YAML 配置，支持 `GOUNO_` 前缀环境变量覆盖：

```yaml
web_server:
    address: 0.0.0.0
    port: 8080
    request_timeout: 10s
    rate_limit_per_minute: 120
```

```bash
# 通过环境变量覆盖端口
GOUNO_WEB_SERVER_PORT=3000 ./bin/my-service web
```

CLI 参数优先级最高，覆盖配置文件和环境变量：

```bash
./bin/my-service web -p 3000 -a 127.0.0.1 -e development
```

## CLI 命令

```bash
my-service web                          # 启动 Web 服务
my-service web -e development           # 使用开发环境配置
my-service web -p 3000 -a 127.0.0.1     # 自定义端口和地址
my-service gen suite user               # 生成 DDD 模块
my-service gen task send_email          # 生成后台任务
my-service migrate up                   # 执行数据库迁移
my-service migrate down 1               # 回滚 1 个迁移
my-service migrate status               # 查看迁移状态
```

## 设计理念

**gouno 是启动器，不是框架。**

| | 启动器（gouno） | 框架 |
|--|--|--|
| **定位** | 处理启动模板代码 | 抽象底层细节 |
| **所有权** | 你拥有代码，随意修改 | 你遵循框架规则 |
| **技术选型** | 你来选 | 框架来选 |
| **扩展方式** | 模板集 | 插件/模块 |

gouno 给你一个标准化的起点，然后让开。真正的定制通过**模板集**完成——你的代码风格、技术选型、团队规范——全部封装为可复用的模板。

```
gouno（核心）         → 标准化：项目结构 + 启动流程 + Web 层 + 响应格式
模板集                → 定制化：代码风格 + 技术选型 + 团队规范
你的业务代码          → 实现：产品逻辑
```

## 相关项目

| 项目 | 说明 |
|------|------|
| [gouno](https://github.com/rushairer/gouno) | 核心库（本仓库） |
| [gouno-cli](https://github.com/rushairer/gouno-cli) | CLI 工具，用于创建项目和管理模板集 |
| [gouno-template](https://github.com/rushairer/gouno-template) | 默认模板集和项目模板 |

## 许可证

MIT License。详见 [LICENSE](LICENSE)。
