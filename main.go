package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"telegramStravaBot/config"
	"telegramStravaBot/data/database"
	"telegramStravaBot/infrastructure"
	"telegramStravaBot/interfaces"
)

func main() {
	infrastructure.LoadEnv()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	configuration, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	_, err = database.Connect(configuration.Database)

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	ya := interfaces.NewYandexWeather(bot)
	ya.Init()

	ui := interfaces.NewTelegramUI()
	telegramRepo := interfaces.TelegramUIRepository{UI: ui, YA: ya}
	telegramRepo.Init(bot)
}
