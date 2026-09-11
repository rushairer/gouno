# Gouno Core Agent Contract

This file defines repository-maintenance invariants for humans and coding agents working on `rushairer/gouno`.

## Product boundary

Gouno is a lightweight project starter/runtime-mechanism library and project-tooling host. It is not the owner of an application's architecture.

Gouno Core may own reusable mechanisms such as protocol parsing, validation, rendering, middleware, response primitives, and other runtime helpers. It must not turn template-specific vocabulary or architectural policy into Core concepts.

In particular, names such as `domain`, `repository`, `service`, `controller`, `handler`, `usecase`, `task`, and `suite` are not Gouno Core concepts. They may appear in examples only when clearly identified as template policy.

## Codegen ownership

The normative Codegen v1 contract is `docs/codegen-template-spec.md`.

Gouno owns the codegen protocol and execution engine. A project template owns whether codegen exists and, when it does, the generator catalog, command shape, arguments, flags, output paths, composition policy, and source templates.

Do not add a concrete generator to Core merely because it is useful in the default `gouno-template`.

`generator.GeneratorCmd` exists as a compatibility proxy. Do not use it as justification to restore a built-in generator catalog.

## Protocol changes

A change to manifest fields, discovery, rendering semantics, path safety, composition, overwrite behavior, or command construction is a public protocol change.

For such changes:

1. assess schema compatibility before editing implementation;
2. update `docs/codegen-template-spec.md` in the same change;
3. add or update contract tests;
4. update README when user-visible behavior changes;
5. record compatibility/user-visible changes in `CHANGELOG.md` under `Unreleased`;
6. consider whether the schema identifier must change rather than silently changing v1 semantics.

Codegen v1 deliberately has no arbitrary shell hooks or executable plugins. Do not introduce arbitrary execution into v1 without an explicit trust model, security review, and versioned protocol decision.

## Cross-repository boundary

Related responsibilities are intentionally split:

- `gouno`: reusable mechanisms and Codegen protocol/runtime;
- `gouno-cli`: project-template bootstrap contract and project creation;
- `gouno-template`: default project/template policy and reference implementation;
- `gouno-doc`: user and template-author guides.

Do not duplicate normative protocol text across repositories. Link to the authoritative contract instead.

## Validation

Follow repository CI and contribution requirements. Public protocol or compatibility changes are not complete until tests, vet, lint, vulnerability checks, and relevant cross-repository integration checks pass.
