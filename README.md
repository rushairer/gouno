# gouno

[中文](./README.zh-CN.md) | [Documentation](https://github.com/rushairer/gouno-doc)

---

A **lightweight Go web project launcher**. It provides reusable startup, HTTP, response, security, and project-tooling primitives while templates decide concrete project structure and development conventions.

gouno is **not a full-stack framework**. It does not prescribe a database, cache, message queue, authentication product, or application architecture.

```text
Gouno owns                          Templates / applications own
├── reusable runtime primitives    ├── project architecture
├── codegen protocol + engine      ├── generator catalog + templates
├── template-safe rendering        ├── database/cache/message queue
├── middleware/security helpers    ├── authentication implementation
└── common project tooling         └── application/domain policy
```

## Supported Go versions

gouno requires Go 1.25 or newer. CI covers the latest security patch of the current and previous stable Go series; applications should keep their patch release current.

## Quick Start

```bash
go install github.com/rushairer/gouno-cli@latest

gouno-cli new my-service -m github.com/you/my-service
cd my-service && go mod tidy && make dev
```

## Template-defined code generation

Code generation is a **template capability**, not a built-in DDD opinion in Gouno.

If the current template ships `.gouno/codegen.yaml`, it can define commands such as:

```bash
gouno gen suite user
gouno gen task send_email
gouno gen controller auth
```

A different template can expose `handler`, `usecase`, `module`, or any other generator names. A template without a codegen manifest can expose no `gen` command at all.

Gouno only implements discovery, manifest validation, CLI construction, argument/flag handling, safe rendering, composition, overwrite policy, and Go formatting. Concrete generator names, parameters, output paths, and source templates belong to the project template.

See [Template Codegen Specification v1](./docs/codegen-template-spec.md).

## Philosophy

**Gouno is a launcher and project-tooling host, not an application framework.** It standardizes mechanisms while leaving policies to templates and applications.

```text
gouno-cli           -> creates a project from any full project template
gouno               -> reusable mechanisms + template-driven tooling protocol
project template    -> bootstrap code + architecture + codegen policy
application         -> product/domain logic
```

A generated project may depend on selected Gouno runtime primitives, but Gouno does not own the project's business architecture. Names such as domain, repository, service, controller, task, and suite are template vocabulary rather than Gouno core concepts.

## Related Projects

| Repository | Description |
|------------|-------------|
| [gouno](https://github.com/rushairer/gouno) | Core library and project-tooling protocol (this repo) |
| [gouno-cli](https://github.com/rushairer/gouno-cli) | Project creation CLI |
| [gouno-template](https://github.com/rushairer/gouno-template) | Default project template and default codegen policy |
| [gouno-doc](https://github.com/rushairer/gouno-doc) | Documentation |

## License

MIT License.

Security reports should follow [SECURITY.md](./SECURITY.md). Contributions and release policy are documented in [CONTRIBUTING.md](./CONTRIBUTING.md) and [RELEASE_CHECKLIST.md](./RELEASE_CHECKLIST.md).
