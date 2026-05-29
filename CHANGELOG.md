# Changelog

## Unreleased

### Added

- English `doc.go`, `errors.go`, expanded README and examples.
- `ClusterAddrs`, connection pool/timeouts, `ConfigDefault()`.
- Injectable client factory hooks; `NewClientWithRetry` alias.
- `config_apply.go` for Options mapping; CI coverage gate ≥ 85%.

### Changed

- `Validate()` applies pool defaults (`PoolSize` 50, `MinIdleConns` 10) and enforces exclusive connection mode.
- Client supports cluster via `redis.Cmdable`; `Redis()` nil for cluster.

### Compatibility

- `NewClient(cfg) (Client, error)` and `Client` interface unchanged for standalone/sentinel callers.
