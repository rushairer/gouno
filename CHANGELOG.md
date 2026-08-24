# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

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
- Removed deprecated `template-set` resolution, `.gouno.yaml` support, and `--template-set` flags from all generator commands to align with `gouno-cli` v1.1.0; renamed file to `generator/template.go` (`generator/template.go`, `generator/*.go`).
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
