package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Init(redisUrl string) error {
	redisOpts, err := redis.ParseURL(redisUrl)

	if err != nil {
		return err
	}

	redisOpts.PoolSize = 500
	redisOpts.MinIdleConns = 100
	redisOpts.DialTimeout = 5 * time.Second
	redisOpts.ReadTimeout = 3 * time.Second
	redisOpts.WriteTimeout = 3 * time.Second
	redisOpts.PoolTimeout = 4 * time.Second

	RedisClient = redis.NewClient(redisOpts)
	return RedisClient.Ping(context.Background()).Err()

}
