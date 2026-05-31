# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

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
