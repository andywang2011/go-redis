package goredis_test

import (
	"context"
	"fmt"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/andywang2011/go-redis"
)

func ExampleNewClient() {
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
	// Output: true
}
