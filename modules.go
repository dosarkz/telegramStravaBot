package main

import (
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
	"telegramStravaBot/config"
	"telegramStravaBot/database/connections"
)

func LoadDependencies() (*tgbotapi.BotAPI, *gorm.DB, *redis.Client) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	configuration, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// establish DB connection
	database, err := connections.Connect(configuration)
	if err != nil {
		panic(err)
	}

	// migration
	m, err := migrate.New(
		"file://"+configuration.BasePath+"/database/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			configuration.Database.User, configuration.Database.Password, configuration.Database.Host,
			configuration.Database.Port, configuration.Database.DB))
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	// establish Redis connection
	redisConn, err := connections.ConnectToRedis(configuration.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	return bot, database, redisConn
}
