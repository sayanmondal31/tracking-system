package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Init(redisUrl string) error {
	redisOpts, err := redis.ParseURL(redisUrl)

	if err != nil {
		return err
	}

	RedisClient = redis.NewClient(redisOpts)
	return RedisClient.Ping(context.Background()).Err()

}
