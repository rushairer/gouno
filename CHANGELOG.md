# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [1.4.0] - 2026-11-04

### Changed
- Require Go 1.26.0 or newer and validate the latest Go 1.26.x and 1.27.x patch releases in CI.
- Upgrade the Go 1.27-compatible golangci-lint toolchain to v2.15.0.
- Update go-pipeline to v2.3.0, YAML to v3.0.6, x/sync to v0.24.0, and x/text to v0.43.0.
- Mark the Codegen v1 specification as the stable normative contract and clarify its single-source-of-truth ownership.

### Security
- Upgrade golang.org/x/net to v0.60.0, x/crypto to v0.58.0, and x/sys to v0.49.0 to remediate reachable HTTP/2 vulnerabilities.

## [1.3.0] - 2026-09-11

### Added
- Add Template Codegen Specification v1 (`gouno.dev/codegen/v1`).
- Add project-aware codegen manifest discovery, validation, dynamic Cobra command construction, declarative render outputs, generator composition, safe path enforcement, and automatic Go formatting.
- Add `LoadProjectCommand` and `AttachProjectCommand`; templates can expose codegen only when they provide `.gouno/codegen.yaml`.

### Changed
- Move concrete generator policy out of Gouno Core. Generator names such as `domain`, `repository`, `service`, `controller`, `task`, and `suite` are now supplied by project templates.
- Keep `GeneratorCmd` only as a compatibility proxy that delegates execution to the current project's codegen manifest.

### Security
- Reject absolute and project-escaping template/output paths before code generation.
- Codegen v1 intentionally does not allow arbitrary shell hooks or executable plugins.

## [1.2.2] - 2026-08-28

### Added
- Add `RequestID` field and `WithRequestID` method on `Response` to support request tracing in JSON API responses.

## [1.2.1] - 2026-08-28

### Fixed
- Reduce the JWKS refresh cooldown so resource servers can accept newly rotated signing keys promptly while retaining single-flight request coalescing.

### Security
- Restrict access-token verification to RS256 instead of accepting other RSA signing algorithms.
- Require a matching `azp` or `client_id` claim whenever resource servers configure a client identifier.

## [1.2.0] - 2026-08-24

### Changed
- Promote `1.2.0-rc.1` to the stable channel without code changes after downstream integration verification.

## [1.2.0-rc.1] - 2026-08-24

### Changed
- Require Go 1.25 or newer and test the latest security patch of the current and previous stable Go series.
- Document the JWT/JWKS, CSRF, and security-header primitives as supported public capabilities.

### Security
- Upgrade `golang.org/x/text` to 0.39.0 and `quic-go` to 0.59.1 to address reachable denial-of-service and memory-exhaustion advisories.

## [1.1.0] - 2026-08-20

### Changed
- Default code generation path for `controller` changed from `controller` to `internal/controller` to align with standard internal package layout (`generator/controller.go`).
- Preset error responses (`InternalServerErrorResponse`, `BadRequestResponse`, etc.) are now supplemented with immutable constructor functions (`NewInternalServerErrorResponse()`, `NewBadRequestResponse()`, etc.) — each call returns a fresh `*Response` instance, eliminating shared mutable state risk. The old package-level variables are preserved as deprecated aliases for backward compatibility (`response.go`).
- Rate limiter now enforces a `maxVisitors` cap (default 10000) on the visitors map — prevents memory exhaustion from large numbers of unique IPs. Use `SetMaxVisitors()` to customize. When the cap is reached, idle visitors are evicted before rejecting new IPs (`middleware/ratelimit.go`).

### Removed
- Removed deprecated `template-set` resolution, `.gouno.yaml` support, and `--template-set` flags from all generator commands to align with `gouno-cli` v1.1.0; renamed file to `generator/template.go` (`generator/*.go`).
- Deprecated package-level mutable response variables (`InternalServerErrorResponse`, `BadRequestResponse`, `UnauthorizedResponse`, `ForbiddenResponse`, `NotFoundResponse`, `MethodNotAllowedResponse`, `RequestTimeoutResponse`, `ConflictResponse`, `GoneResponse`) — use the corresponding `New*Response()` factory functions instead (`response.go`).

## [1.0.0] - 2026-05-31

### Added

- Unified JSON response format (`Response`, `NewSuccessResponse`, `NewErrorResponse`).
- 9 preset error responses (`ErrBadRequest`, `ErrUnauthorized`, etc.).
- Rate limiter middleware with sliding window algorithm (`RateLimitMiddleware`, `IPRateLimitMiddleware`).
- Code generator framework (`GeneratorCmd`) supporting `domain`, `repository`, `service`, `controller`, `task`, `suite` subcommands.
- Task pipeline abstraction (`Task` interface, `NewTaskPipeline`).
- String utilities: `ToCamelCase`, `ToSnakeCase`.
- Full godoc documentation for all exported symbols.
- Comprehensive unit tests for all packages including `task/`.
