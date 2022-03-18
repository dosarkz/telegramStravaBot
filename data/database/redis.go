package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"telegramStravaBot/config"
)

var ctx = context.Background()

func ConnectToRedis(configuration *config.Redis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configuration.Host, configuration.Port),
		Password: "",
		DB:       0,
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		return nil, err
	}

	_, err = rdb.Get(ctx, "key").Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
