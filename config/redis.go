package config

import (
	"github.com/redis/go-redis/v9"
)

func RedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb, nil
}
