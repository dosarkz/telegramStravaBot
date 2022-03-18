package config

import (
	"os"
)

// Config is a struct that contains configuration variables
type Config struct {
	Environment string
	Port        string
	Database    *Database
	BasePath    string
	BotName     string
	Redis       *Redis
}

// Database is a struct that contains DB's configuration variables
type Database struct {
	Host     string
	Port     string
	User     string
	DB       string
	Password string
	Timezone string
}

type Redis struct {
	Host      string
	Port      string
	Password  string
	CacheTime string
}

// NewConfig creates a new Config struct
func NewConfig() (*Config, error) {
	port := os.Getenv("PORT")
	// set default PORT if missing
	if port == "" {
		port = "3000"
	}
	return &Config{
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
