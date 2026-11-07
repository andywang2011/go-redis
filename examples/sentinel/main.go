// Example Sentinel failover client (requires running Sentinel + Redis).
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andywang2011/go-redis"
)

func main() {
	cfg := redis.Config{
		MasterName:    "mymaster",
		SentinelAddrs: []string{"127.0.0.1:26379", "127.0.0.1:26380"},
		Password:      "",
	}
	cli, err := redis.NewClient(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	if err := cli.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("sentinel ok")
}
