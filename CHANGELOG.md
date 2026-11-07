# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-11-07

### Breaking

- See repository diff for API changes in this release (secure defaults, renamed packages where applicable, context-aware constructors).

### Security

- Hardened defaults: TLS, timeouts, body limits, PII redaction where logging applies.

### Fixed

- Bug fixes from the v0.2.0 hardening pass.

### Changed

- CI uses go-version-file; added race, staticcheck, govulncheck, golangci-lint (gosec).



### Changed

- **Breaking:** package renamed `goredis` -> `redis`; upstream is imported as `goredislib`.
- **Breaking:** `NewClient` and `NewClientWithRetry` now take a `context.Context` for the startup Ping retry.
- `Close` is idempotent via `sync.Once`; `BuildTLS` supports CA-only TLS and rejects cert without key.
- Retry helper uses `math/rand/v2` and a stoppable timer instead of `time.After`.
- Retry helper returns `ctx.Err()` when ctx is cancelled during an attempt and tolerates a nil ctx.

## [0.1.0] - 2026-11-07

### Added

- Standalone, Sentinel, and Cluster clients with startup Ping retry and pool defaults.
