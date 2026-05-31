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

## Template Sets

Customize what `gouno gen` produces. Different teams, different code styles — all without touching gouno's source.

```bash
gouno-cli template install gorm https://github.com/myorg/gouno-template-gorm
gouno-cli new order-service --template-set gorm -m github.com/myorg/order-service
```

[Create your own template set →](https://github.com/rushairer/gouno-doc/blob/main/template-sets.md)

## Documentation

| Guide | Description |
|-------|-------------|
| [Getting Started](https://github.com/rushairer/gouno-doc/blob/main/getting-started.md) | Install, create project, run |
| [Code Generation](https://github.com/rushairer/gouno-doc/blob/main/code-generation.md) | Generate DDD modules |
| [Template Sets](https://github.com/rushairer/gouno-doc/blob/main/template-sets.md) | Create and share custom templates |
| [Configuration](https://github.com/rushairer/gouno-doc/blob/main/configuration.md) | Multi-environment YAML config |
| [Middleware](https://github.com/rushairer/gouno-doc/blob/main/middleware.md) | Built-in and custom middleware |

## Philosophy

**gouno is a launcher, not a framework.** It gives you a standardized starting point, then gets out of the way. The real customization happens in **template sets** — your team's code style, your tech stack choices, your patterns — all captured as reusable templates.

```
gouno (core)          → Standardizes: project structure + startup + web layer + response
template sets         → Customizes: code style + tech stack + team conventions
your business code    → Implements: actual product logic
```

## Related Projects

| Repository | Description |
|------------|-------------|
| [gouno](https://github.com/rushairer/gouno) | Core library (this repo) |
| [gouno-cli](https://github.com/rushairer/gouno-cli) | CLI tool |
| [gouno-template](https://github.com/rushairer/gouno-template) | Default template set |
| [gouno-doc](https://github.com/rushairer/gouno-doc) | Documentation |

## License

MIT License.
