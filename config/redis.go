package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func RedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb, rdb.Ping(ctx).Err()
}
