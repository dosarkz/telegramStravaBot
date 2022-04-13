package main

import (
	"log"
	"telegramStravaBot/domain"
	"telegramStravaBot/domain/bot"
	"telegramStravaBot/domain/yandex"
)

func main() {
	tgBot, db, _ := LoadDependencies()
	log.Printf("Authorized on account %s", tgBot.Self.UserName)

	ya := yandex.NewYandexWeather(tgBot)
	ya.Init()

	ui := bot.NewUIService(bot.UIActionService{Bot: tgBot, YA: ya}, domain.New(db))
	ui.Run()
}
