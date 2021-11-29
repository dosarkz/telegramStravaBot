package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"telegramStravaBot/config"
	db "telegramStravaBot/data/database"
	userStore "telegramStravaBot/data/users"
	workoutStore "telegramStravaBot/data/workouts"
	"telegramStravaBot/domain/users"
	"telegramStravaBot/domain/workouts"
	"telegramStravaBot/infrastructure"
	"telegramStravaBot/interfaces"
	router "telegramStravaBot/router/http"
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

	// establish DB connection
	database, err := db.Connect(configuration)
	if err != nil {
		panic(err)
	}

	// initialize repos and services using DI
	userRepo := userStore.New(database)
	userService := users.NewService(userRepo)

	workoutRepo := workoutStore.New(database)
	workoutService := workouts.NewService(workoutRepo)

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	ya := interfaces.NewYandexWeather(bot)
	ya.Init()

	httpRouter := router.NewHTTPHandler(workoutService)
	err = http.ListenAndServe(":"+configuration.Port, httpRouter)
	fmt.Printf("Connect port %s", configuration.Port)
	if err != nil {
		panic(err)
	}

	ui := interfaces.NewTelegramUI()
	telegramRepo := interfaces.TelegramUIRepository{UI: ui, YA: ya, User: userService, Workout: workoutService,
		Bot: bot}
	telegramRepo.Init()

}
