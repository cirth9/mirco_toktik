package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "47.115.217.189:6379",
		Password: "qq31415926535--",
		DB:       0,
	})
	fmt.Println(len(client.SMembers(context.Background(), "1237198273281").Val()) == 0)
	fmt.Println(client.Get(context.Background(), "1237198273281").Result())
	if client.Get(context.Background(), "1237198273281").Err() != nil {
		fmt.Println(1)
	}
}
