// Package redis provides a thin, testable wrapper around Redis clients (github.com/redis/go-redis/v9).
//
// Design constraints (for use in application frameworks):
//   - No viper, no global singletons, no log.Fatal inside the library.
//   - Client interface embeds goredislib.Cmdable for mock-friendly tests (miniredis).
//   - One Client per process; share it across goroutines (connection pool is thread-safe).
//   - Call Close on shutdown.
//
// High availability: Sentinel failover (MasterName + SentinelAddrs) or Redis Cluster (ClusterAddrs).
// Startup connect uses Ping with exponential-backoff retry (RetryConfig).
package redis
