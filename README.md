# go-redis

Thin, testable Redis client wrapper for Go (`github.com/redis/go-redis/v9`). Package name `goredis`.

```bash
go get github.com/andywang2011/go-redis@latest
```

## Design

- No viper, no global singletons, no `log.Fatal` inside the library.
- `Client` embeds `redis.Cmdable` (mock with miniredis in tests).
- Startup **Ping** uses `RetryConfig`; command-level retry uses driver `MaxRetries`.
- Supports standalone, Sentinel failover, and Redis Cluster.

## Quick start

```go
import goredis "github.com/andywang2011/go-redis"

cfg := goredis.Config{Address: "127.0.0.1:6379"}
cli, err := goredis.NewClient(cfg)
if err != nil { /* handle */ }
defer cli.Close()

_ = cli.Set(ctx, "key", "value", 0).Err()
```

Sentinel:

```go
cfg := goredis.Config{
    MasterName: "mymaster",
    SentinelAddrs: []string{"127.0.0.1:26379"},
}
```

Cluster:

```go
cfg := goredis.Config{ClusterAddrs: []string{"127.0.0.1:7000", "127.0.0.1:7001"}}
```

## Configuration

| Field | Default / notes |
|-------|-----------------|
| `Address` / `MasterName`+`SentinelAddrs` / `ClusterAddrs` | Exactly one mode |
| `PoolSize` | `50` |
| `MinIdleConns` | `10` |
| `MaxRetries` | `2` (driver command retry) |
| `Protocol` | `2` (RESP3) |
| `RetryConfig` | Startup Ping retry; `DefaultRetryConfig()` when unset |
| Timeouts | `DialTimeout`, `ReadTimeout`, `WriteTimeout`, `PoolTimeout`, `ConnMaxIdleTime`, `ConnMaxLifetime` |

## Production

- **One** `Client` per process; share across goroutines.
- Use Sentinel or Cluster for HA; tune pool size for concurrent load.
- `Redis()` returns `*redis.Client` only for standalone/sentinel; cluster callers use `Cmdable` only.
- Call `Close()` on shutdown.

## Limitations

- No built-in circuit breaker (use app-layer resilience, e.g. demo-api `pkg/goredis`).
- Cluster mode: `Redis()` is nil.

## 中文摘要

Redis 轻量封装：单节点 / Sentinel / Cluster、连接池与超时可调、启动 Ping 重试。`Client` 接口兼容 `Cmdable`，便于 miniredis 单测与 Wire 注入。无 viper/全局变量。
