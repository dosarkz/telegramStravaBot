package database

import (
	"fmt"
	"github.com/go-redis/redis"
	"telegramStravaBot/config"
)

func ConnectToRedis(configuration *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configuration.Host, configuration.Port),
		Password: configuration.Password,
		DB:       0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}
