package config

import (
	"os"
)

// App is a struct that contains configuration variables
type App struct {
	Environment string
	Port        string
	Database    *Database
	BasePath    string
	BotName     string
	Redis       *Redis
}

// NewConfig creates a new Config struct
func NewConfig() (*App, error) {
	port := os.Getenv("PORT")
	// set default PORT if missing
	if port == "" {
		port = "3000"
	}
	return &App{
		Environment: os.Getenv("ENV"),
		Port:        port,
		BasePath:    os.Getenv("BASE_PATH"),
		BotName:     os.Getenv("BOT_NAME"),
		Database: &Database{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			User:     os.Getenv("DATABASE_USER"),
			DB:       os.Getenv("DATABASE_DB"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Timezone: os.Getenv("DATABASE_TIMEZONE"),
		},
		Redis: &Redis{
			Host:      os.Getenv("REDIS_HOST"),
			Port:      os.Getenv("REDIS_PORT"),
			Password:  os.Getenv("REDIS_PASSWORD"),
			CacheTime: os.Getenv("CACHE_TIME"),
		},
	}, nil
}
