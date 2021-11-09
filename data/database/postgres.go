package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"telegramStravaBot/config"
)

func Connect(config *config.Database) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
			config.Host, config.Port, config.DB, config.User, config.Password),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
