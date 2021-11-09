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
		Database: &Database{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			User:     os.Getenv("DATABASE_USER"),
			DB:       os.Getenv("DATABASE_DB"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Timezone: os.Getenv("DATABASE_TIMEZONE"),
		},
	}, nil
}
