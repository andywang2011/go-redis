# go-redis

Thin, testable Redis client wrapper (package `goredis`).

## Install

```bash
go get github.com/andywang2011/go-redis@v0.1.0
```

## Quick start

```go
import (
	"context"
	"fmt"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/andywang2011/go-redis"
)

mr, err := miniredis.Run()
if err != nil {
	fmt.Println("redis error")
	return
}
defer mr.Close()

cfg := goredis.Config{Address: mr.Addr()}
cli, err := goredis.NewClient(cfg)
if err != nil {
	fmt.Println("client error")
	return
}
defer cli.Close()

err = cli.Set(context.Background(), "key", "value", 0).Err()
fmt.Println(err == nil)
```

## Examples

```bash
export REDIS_ADDR=127.0.0.1:6379
cd go-redis/examples/basic && go run .
cd go-redis/examples/sentinel && go run .
```

## License

MIT — see [LICENSE](LICENSE).
