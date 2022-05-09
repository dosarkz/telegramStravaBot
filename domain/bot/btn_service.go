package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

type UIButtonInt interface {
	StravaInlineButton() tgbotapi.InlineKeyboardButton
	InstaInlineButton() tgbotapi.InlineKeyboardButton
	Participate(text string, callback string) tgbotapi.InlineKeyboardButton
	MetronomeInlineButton() tgbotapi.InlineKeyboardButton
	BotanInlineButton() tgbotapi.InlineKeyboardButton
}

type UIButtonService struct{}

func (t *UIButtonService) StravaInlineButton() tgbotapi.InlineKeyboardButton {
	stravaButtonData := os.Getenv("METRONOME_STRAVA_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Подписаться в Strava", URL: &stravaButtonData}
}

func (t *UIButtonService) AlmatyInlineButton() tgbotapi.InlineKeyboardButton {
	almatyButtonData := os.Getenv("ALMATY_GROUP_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Алматинское сообщество бега", URL: &almatyButtonData}
}

func (t *UIButtonService) CompleteWorkoutButton() tgbotapi.InlineKeyboardButton {
	str := "workout_complete"
	return tgbotapi.InlineKeyboardButton{Text: "❌ Завершить", CallbackData: &str}
}

func (t *UIButtonService) UpdateHeroButton() tgbotapi.InlineKeyboardButton {
	str := "update_hero"
	return tgbotapi.InlineKeyboardButton{Text: "🔃 Обновить", CallbackData: &str}
}

func (t *UIButtonService) UpdateMetroButton() tgbotapi.InlineKeyboardButton {
	str := "update_metro"
	return tgbotapi.InlineKeyboardButton{Text: "🔃 Обновить", CallbackData: &str}
}

func (t *UIButtonService) UpdatingHeroButton() tgbotapi.InlineKeyboardButton {
	str := "update_hero"
	return tgbotapi.InlineKeyboardButton{Text: "🏃 Разминаюсь и бегу, пожалуйста ждите...", CallbackData: &str}
}

func (t *UIButtonService) InstaInlineButton() tgbotapi.InlineKeyboardButton {
	instagramButtonData := os.Getenv("METRONOME_INSTA_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Подписаться в Instagram", URL: &instagramButtonData}
}

func (t *UIButtonService) MetronomeInlineButton() tgbotapi.InlineKeyboardButton {
	metronomeButtonData := os.Getenv("METRONOME_TELEGRAM_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Сообщество бега в Триатлон парке (Metronome)", URL: &metronomeButtonData}
}

func (t *UIButtonService) BotanInlineButton() tgbotapi.InlineKeyboardButton {
	botanButtonData := os.Getenv("BOTAN_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Сообщество бега в Ботаническом саду (Ботаны)", URL: &botanButtonData}
}

func (t *UIButtonService) CentralInlineButton() tgbotapi.InlineKeyboardButton {
	centralButtonData := os.Getenv("CENTRAL_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Сообщество бега в Центральном парке (whatsapp группа)", URL: &centralButtonData}
}

func (t *UIButtonService) AlmatyRunningGroups() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{t.AlmatyInlineButton()})
}

func (t *UIButtonService) Participate(text string, callback string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{Text: text, CallbackData: &callback}
}
