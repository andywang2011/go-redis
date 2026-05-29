// Example standalone Redis client: set, get, close.
package main

import (
	"context"
	"fmt"
	"log"

	goredis "github.com/andywang2011/go-redis"
)

func main() {
	cli, err := goredis.NewClient(goredis.Config{Address: "127.0.0.1:6379"})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()
	if err := cli.Set(ctx, "hello", "world", 0).Err(); err != nil {
		log.Fatal(err)
	}
	val, err := cli.Get(ctx, "hello").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
}
