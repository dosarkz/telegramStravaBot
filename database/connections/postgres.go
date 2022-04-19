package connections

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Import GORM postgres dialect for its side effects, according to GORM docs.
	"telegramStravaBot/config"
)

func Connect(config *config.App) (*gorm.DB, error) {

	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.DB, config.Database.User, config.Database.Password)
	db, err := gorm.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}
