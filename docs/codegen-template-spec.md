# Gouno Template Codegen Specification v1

> Status: **Stable / Normative**  
> Schema: `gouno.dev/codegen/v1`  
> Introduced: `gouno v1.3.0`  
> Authority: `rushairer/gouno`

This document is the normative contract for Gouno Codegen v1. User guides and reference templates may explain or demonstrate this contract, but must not redefine it. If implementation and this specification disagree, treat the mismatch as a compatibility defect that must be resolved explicitly rather than silently changing v1 semantics.

Gouno owns the code-generation protocol and runtime. A project template owns the generator catalog, CLI shape, paths, and source templates.

This keeps Gouno project-aware without making Gouno opinionated about DDD, Clean Architecture, handlers, repositories, or any other application structure.

## Capability discovery

A generated project supports code generation only when it contains:

```text
.gouno/codegen.yaml
```

`generator.AttachProjectCommand(root, startDir)` searches the current directory and its parents for that file. If it is absent, no codegen command is attached to the project CLI.

Templates that do not want code generation simply omit the manifest and do not need to ship generator templates.

## Manifest

The v1 schema identifier is:

```yaml
schema: gouno.dev/codegen/v1
```

Example:

```yaml
schema: gouno.dev/codegen/v1
command:
  use: gen
  short: Generate project code
  aliases: [generator]

generators:
  - name: service
    aliases: [s]
    short: Generate service
    args:
      - name: name
        required: true
    flags:
      - name: path
        shorthand: p
        type: string
        default: internal/service
        description: path to service
      - name: force
        shorthand: f
        type: bool
        default: false
        description: force overwrite
    outputs:
      - template: .gouno/codegen/service.tmpl
        path: '{{ flag "path" }}/{{ arg "name" }}.go'
```

The resulting command is defined by the template, not by Gouno:

```text
gen service [name] --path internal/service --force
# alias: generator service [name]
```

## Generator fields

- `name`: command name. Must be unique.
- `short`: Cobra short description.
- `aliases`: optional command aliases.
- `args`: ordered positional arguments. Required arguments must precede optional arguments.
- `flags`: v1 supports `string` and `bool` flags.
- `outputs`: files rendered by this generator.
- `compose`: names of other generators to invoke in order.

A generator must provide at least one output or composed generator.

### Reserved `force` flag

If a generator declares a boolean flag named `force`, Gouno uses it as the overwrite policy for rendered files. Without `force=true`, existing output files are skipped.

A composition generator can declare the same flag; matching flags are inherited by composed generators. This is how a `suite` generator can propagate `--force` without Gouno knowing what a suite is.

## Template expressions

Output paths and template contents use Go `text/template` with these functions:

- `arg "name"`: read a positional argument.
- `flag "path"`: read a flag value as text.
- `camel value`: convert to PascalCase.
- `snake value`: convert to snake_case.
- `kebab value`: convert to kebab-case.
- `lower value`: lowercase.
- `upper value`: uppercase.

The raw maps are also available as `.Args` and `.Flags`.

## Composition

Composition is declarative in v1:

```yaml
- name: suite
  short: Generate domain, repository and service
  args:
    - name: name
      required: true
  flags:
    - name: force
      shorthand: f
      type: bool
      default: false
      description: force overwrite
  compose:
    - domain
    - repository
    - service
```

Composed generators receive arguments with matching names and flags with matching names. Cycles and missing composed generators are rejected when the manifest loads.

## Safety and determinism

Codegen v1 deliberately does not support arbitrary shell hooks or executable plugins.

The engine:

- rejects output and template paths that escape the project root;
- refuses absolute paths;
- formats `.go` output with `go/format` before writing;
- skips existing files unless the template explicitly exposes `force` and the user enables it;
- validates the manifest before constructing the command tree.

External tools, AST patching, and executable hooks may be considered in a future schema only with an explicit trust model.

## Template authoring contract

A template that supports codegen SHOULD ship the manifest and referenced templates with the generated project so generator behavior remains pinned to the project's original template version.

A template MUST NOT depend on the latest remote template at generation time. Updating a project's codegen policy should be an explicit template migration, not an implicit network lookup.

Recommended layout:

```text
.gouno/
  codegen.yaml
  codegen/
    controller.tmpl
    domain.tmpl
    repository.tmpl
    service.tmpl
    task.tmpl
```

The `.gouno/codegen/` subtree is runtime template data. Project-bootstrap renderers should copy it verbatim so its own Go-template expressions remain intact.

Gouno does not assign meaning to names such as `domain`, `service`, `handler`, or `usecase`. Those are template policy.
