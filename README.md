# gouno

[中文](./README.zh-CN.md)

---

A **lightweight Go web project launcher**. Scaffolds project structure, startup flow, web layer, and response format — the boilerplate you'd rewrite for every microservice — so you can focus on business logic from the first line.

gouno is **not a framework**. It doesn't bundle a database driver, cache client, or message queue. Your tech stack is your choice. gouno handles the rest.

```
What gouno does                    What gouno doesn't do
├── Project structure (DDD)        ├── Database (pgx? gorm? ent?)
├── CLI + config (Cobra + Viper)   ├── Cache (redis? memcached?)
├── Web engine (Gin + middleware)   ├── Messaging (kafka? rabbitmq?)
├── Response format (unified JSON)  └── Auth (JWT? OAuth? session?)
└── Code generator + templates
```

## Quick Start

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

## Code Generation

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

Aliases for faster typing:

```bash
gouno gen d order   # domain
gouno gen r order   # repository
gouno gen s order   # service
gouno gen c order   # controller
gouno gen t order   # task
```

Options:

```bash
gouno gen suite user --force          # Overwrite existing files
gouno gen suite user --path ./custom  # Custom output directory
```

## Template Sets

Different teams, different code styles. Template sets let you customize what `gouno gen` produces — without touching gouno's source code.

### Install and Use

```bash
# Install a community or company template set
gouno-cli template install gorm https://github.com/myorg/gouno-template-gorm

# Use it when creating a project
gouno-cli new order-service --template-set gorm -m github.com/myorg/order-service

# Now gen suite uses your gorm templates
cd order-service
gouno gen suite user   # → generates gorm-style domain/repository/service
```

### Manage Template Sets

```bash
gouno-cli template list                        # List installed sets
gouno-cli template install <name> <git-url>    # Install from git
gouno-cli template install <name> <local-path> # Install from local directory
gouno-cli template remove <name>               # Remove a set
gouno-cli template install <name> <url> --force # Overwrite existing
```

### How Template Search Works

```
Search priority:
1. --template-set flag (command level)
2. .gouno.yaml in project root (project level)
3. Built-in default templates (fallback)

Installed location:
~/.gouno/templates/
├── default/      ← Ships with gouno-template
├── gorm/         ← Custom set
└── my-company/   ← Company set
```

### Create Your Own Template Set

A template set is a directory with `.tmpl` files:

```
my-template-set/
├── domain.tmpl
├── repository.tmpl
├── service.tmpl
├── controller.tmpl
└── task.tmpl
```

Each `.tmpl` file uses `%s` as the struct name placeholder (5 occurrences):

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

Distribute your template set as a git repository. Others install it with one command:

```bash
gouno-cli template install my-set https://github.com/myorg/my-template-set
```

## Project Structure

```
my-service/
├── cmd/
│   ├── main.go              ← Entry point
│   └── gouno/
│       ├── root.go          ← Cobra root command
│       └── web.go           ← Web server startup
├── config/
│   ├── development.yaml     ← Dev config
│   ├── test.yaml            ← Test config
│   ├── production.yaml      ← Production config (env vars for secrets)
│   ├── config.go            ← Config struct definitions
│   └── config_manager.go    ← Thread-safe config loader
├── controller/              ← HTTP handlers
├── internal/
│   ├── domain/              ← Entities
│   ├── repository/          ← Data access interfaces
│   ├── service/             ← Business logic
│   └── task/                ← Background tasks
├── middleware/              ← Custom middleware
├── router/                  ← Route registration
├── utility/                 ← Shared utilities
├── .gouno.yaml              ← Template set config (optional)
├── Makefile                 ← Build/run/dev shortcuts
├── .air.toml                ← Hot-reload config
└── go.mod
```

## API Response Format

All responses use a unified JSON structure. The `code` field matches the HTTP status code:

```json
// Success (HTTP 200)
{"code": 200, "message": "success", "data": {...}}

// Error (HTTP 400/401/403/404/408/429/500)
{"code": 400, "message": "bad request", "data": null}
{"code": 404, "message": "not found", "data": null}
{"code": 500, "message": "internal server error", "data": null}
```

Usage in handlers:

```go
// Success
ctx.JSON(http.StatusOK, gouno.NewSuccessResponse(data))

// Error (pre-defined)
ctx.JSON(http.StatusBadRequest, gouno.BadRequestResponse)
ctx.JSON(http.StatusNotFound, gouno.NotFoundResponse)
ctx.JSON(http.StatusInternalServerError, gouno.InternalServerErrorResponse)

// Error (custom message)
ctx.JSON(http.StatusBadRequest, gouno.NewErrorResponse(http.StatusBadRequest, "email already exists"))
```

## Middleware

Built-in middleware chain (applied in order):

| Middleware | Purpose | HTTP Status on Error |
|------------|---------|---------------------|
| Logger | Request/response logging | — |
| Recovery | Panic recovery | 500 |
| Timeout | Request timeout | 408 |
| RateLimit | IP-based sliding window rate limiting | 429 |

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

## Configuration

Multi-environment YAML config with Viper. Environment variables override via `GOUNO_` prefix:

```yaml
web_server:
    address: 0.0.0.0
    port: 8080
    request_timeout: 10s
    rate_limit_per_minute: 120
```

```bash
# Override port via env var
GOUNO_WEB_SERVER_PORT=3000 ./bin/my-service web
```

CLI flags override both config files and env vars:

```bash
./bin/my-service web -p 3000 -a 127.0.0.1 -e development
```

## CLI Commands

```bash
my-service web                          # Start web server
my-service web -e development           # Use development config
my-service web -p 3000 -a 127.0.0.1     # Custom port and address
my-service gen suite user               # Generate DDD module
my-service gen task send_email          # Generate task
my-service migrate up                   # Run database migrations
my-service migrate down 1               # Rollback 1 migration
my-service migrate status               # Show migration status
```

## Philosophy

**gouno is a launcher, not a framework.**

| | Launcher (gouno) | Framework |
|--|--|--|
| **Role** | Handles startup boilerplate | Abstracts away details |
| **Ownership** | You own the code, modify anything | You follow its rules |
| **Tech stack** | Your choice | Framework's choice |
| **Extension** | Template sets | Plugins/modules |

gouno gives you a standardized starting point, then gets out of the way. The real customization happens in **template sets** — your team's code style, your tech stack choices, your patterns — all captured as reusable templates.

```
gouno (core)          → Standardizes: project structure + startup + web layer + response
template sets         → Customizes: code style + tech stack + team conventions
your business code    → Implements: actual product logic
```

## Related Projects

| Project | Description |
|---------|-------------|
| [gouno](https://github.com/rushairer/gouno) | Core library (this repo) |
| [gouno-cli](https://github.com/rushairer/gouno-cli) | CLI tool for creating projects and managing templates |
| [gouno-template](https://github.com/rushairer/gouno-template) | Default template set and project template |

## License

MIT License. See [LICENSE](LICENSE) for details.
