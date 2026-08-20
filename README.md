# gouno

[中文](./README.zh-CN.md) | [Documentation](https://github.com/rushairer/gouno-doc)

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

