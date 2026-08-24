# gouno

[中文](./README.zh-CN.md) | [Documentation](https://github.com/rushairer/gouno-doc)

---

A **lightweight Go web project launcher**. Scaffolds project structure, startup flow, web layer, and response format — the boilerplate you'd rewrite for every microservice — so you can focus on business logic from the first line.

gouno is **not a full-stack framework**. It doesn't bundle a database driver,
cache client, or message queue. It provides reusable project, HTTP, response,
JWT/JWKS verification, CSRF, and security-header primitives while applications
retain responsibility for their authorization policy and infrastructure choices.

```
What gouno does                    What gouno doesn't do
├── Project structure (DDD)        ├── Database (pgx? gorm? ent?)
├── CLI + config (Cobra + Viper)   ├── Cache (redis? memcached?)
├── Web engine (Gin + middleware)   ├── Messaging (kafka? rabbitmq?)
├── Response format (unified JSON)  ├── Identity provider / login UI
├── Auth + middleware primitives    └── Application authorization policy
└── Code generator + templates
```

## Supported Go versions

gouno v1.2.0 requires Go 1.25 or newer. Each new minor release may raise the
minimum Go version. CI covers the latest security patch of the current and
previous stable Go series; applications should keep their patch release current.

## Quick Start

```bash
# Install
go install github.com/rushairer/gouno-cli@latest

# Create, build, run
gouno-cli new my-service -m github.com/you/my-service
cd my-service && go mod tidy && make dev
# → http://localhost:8080
```

## Code Generation

```bash
gouno gen suite user   # → domain + repository + service
gouno gen task send_email
gouno gen controller auth
```

[Full guide →](https://github.com/rushairer/gouno-doc/blob/main/code-generation.md)

## Documentation

| Guide | Description |
|-------|-------------|
| [Getting Started](https://github.com/rushairer/gouno-doc/blob/main/getting-started.md) | Install, create project, run |
| [Code Generation](https://github.com/rushairer/gouno-doc/blob/main/code-generation.md) | Generate DDD modules |
| [Configuration](https://github.com/rushairer/gouno-doc/blob/main/configuration.md) | Multi-environment YAML config |
| [Middleware](https://github.com/rushairer/gouno-doc/blob/main/middleware.md) | Built-in and custom middleware |

## Philosophy

**gouno is a launcher, not a framework.** It gives you a standardized starting point, then gets out of the way.

```
gouno (core)          → Standardizes: project structure + startup + web layer + response
gouno-template        → Project skeleton: DDD architecture + config + gin + viper
your business code    → Implements: actual product logic
```

## Related Projects

| Repository | Description |
|------------|-------------|
| [gouno](https://github.com/rushairer/gouno) | Core library (this repo) |
| [gouno-cli](https://github.com/rushairer/gouno-cli) | CLI tool |
| [gouno-template](https://github.com/rushairer/gouno-template) | Default project template |
| [gouno-doc](https://github.com/rushairer/gouno-doc) | Documentation |

## License

MIT License.

Security reports should follow [SECURITY.md](./SECURITY.md). Contributions and
release policy are documented in [CONTRIBUTING.md](./CONTRIBUTING.md) and
[RELEASE_CHECKLIST.md](./RELEASE_CHECKLIST.md).
