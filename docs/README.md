# Gouno Core documentation index

Gouno Core is an **architecture-neutral Go runtime/tooling library**. The current code and [root README](../README.md) determine the implemented features. User-facing guides live in [gouno-doc](https://github.com/rushairer/gouno-doc), while the actual Codegen specification is owned here.

## Normative contract

- [Template Codegen Specification v1](codegen-template-spec.md) — maintained schema and protocol for `gouno.dev/codegen/v1`; generator names and project structure remain project/template policy.
- The `gouno-cli` [Project Template Contract](https://github.com/rushairer/gouno-cli/blob/main/docs/project-template-contract.md) governs bootstrap, not Codegen runtime.
- The [default Gouno template](https://github.com/rushairer/gouno-template) is one reference policy, not Core architecture.

## Point-in-time release and security evidence

- [2026-11 ecosystem release audit](quality/ecosystem-release-audit-2026-11.md), [machine-readable audit](quality/ecosystem-release-evidence-2026-11.json) and Gosso v1.7.0 verification receipts record **historical versions at the time measured**. Later formal releases exist; these JSON files are historical evidence, not live compatibility matrices.
- For current published relationships consult the maintained [gouno-doc compatibility matrix](https://github.com/rushairer/gouno-doc/blob/main/standards/compatibility.md) and the exact releases of each project.
- Do not delete original signed-image, checksum or audit receipts simply because a newer artifact superseded them.

All hand-maintained Gouno Core documents belong under `docs/`; there is no `doc/` duplicate. Preserve chronological CHANGELOG history and update normative contracts only alongside matching implementation and tests.
