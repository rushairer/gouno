# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Changed
- Preset error responses (`InternalServerErrorResponse`, `BadRequestResponse`, etc.) are now supplemented with immutable constructor functions (`NewInternalServerErrorResponse()`, `NewBadRequestResponse()`, etc.) — each call returns a fresh `*Response` instance, eliminating shared mutable state risk. The old package-level variables are preserved as deprecated aliases for backward compatibility (`response.go`).
- Rate limiter now enforces a `maxVisitors` cap (default 10000) on the visitors map — prevents memory exhaustion from large numbers of unique IPs. Use `SetMaxVisitors()` to customize. When the cap is reached, idle visitors are evicted before rejecting new IPs (`middleware/ratelimit.go`).

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
