package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegramStravaBot/config"
	"telegramStravaBot/domain/yandex"
)

type UIActionService struct {
	Bot    *tgbotapi.BotAPI
	Config *config.App
	YA     yandex.YandexWeather
}
