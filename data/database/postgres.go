package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Import GORM postgres dialect for its side effects, according to GORM docs.
	"log"
	"telegramStravaBot/config"
)

func Connect(config *config.Config) (*gorm.DB, error) {

	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.DB, config.Database.User, config.Database.Password)
	db, err := gorm.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	// migration
	m, err := migrate.New(
		"file://"+config.BasePath+"/data/database/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.Database.User, config.Database.Password, config.Database.Host,
			config.Database.Port, config.Database.DB))
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	return db, nil
}
