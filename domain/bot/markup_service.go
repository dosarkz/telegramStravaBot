package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type MenuInt interface {
	MainMenu() tgbotapi.ReplyKeyboardMarkup
	HideMenu() tgbotapi.ReplyKeyboardRemove
	MetroInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup
	MarathonInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup
	AppointmentKeyboardMarkup(workoutId int) tgbotapi.InlineKeyboardMarkup
	AppointmentDoneKeyboardMarkup(workoutId int) tgbotapi.InlineKeyboardMarkup
}

type UIMenuService struct {
	Button *UIButtonService
}

func (m *UIMenuService) MainMenu() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🏃 Клуб Любителей Бега MaratHON"),
		),
		tgbotapi.NewKeyboardButtonRow(
			//tgbotapi.NewKeyboardButton("💥 Герой дня"),
			tgbotapi.NewKeyboardButton("⚡ Рейтинг метронома"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("😊 Разминка"),
			tgbotapi.NewKeyboardButton("☂ Погода"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✅ Записаться"),
			tgbotapi.NewKeyboardButton("➕Добавить тренировку"),
		),
	)
	return keyboard
}

func (m *UIMenuService) HideMenu() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
}

func (m *UIMenuService) MetroInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{m.Button.StravaInlineButton(),
		m.Button.InstaInlineButton()})
}

func (m *UIMenuService) MarathonInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{m.Button.MetronomeInlineButton()},
		[]tgbotapi.InlineKeyboardButton{m.Button.BotanInlineButton()},
		[]tgbotapi.InlineKeyboardButton{m.Button.CentralInlineButton()},
	)
}

func (m *UIMenuService) CreateWorkoutKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{m.Button.CompleteWorkoutButton()},
	)
}

func (m *UIMenuService) HeroUpdateButtonKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{m.Button.UpdateHeroButton()},
	)
}

func (m *UIMenuService) MetroUpdateButtonKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{m.Button.UpdateMetroButton()},
	)
}

func (m *UIMenuService) HeroUpdatingButtonKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{m.Button.UpdatingHeroButton()},
	)
}

func (m *UIMenuService) AppointmentKeyboardMarkup(workoutId int) tgbotapi.InlineKeyboardMarkup {
	callbackData := "appointment_" + strconv.Itoa(workoutId)
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{m.Button.Participate("✅Приму участие", callbackData)},
	)
}

func (m *UIMenuService) AppointmentDoneKeyboardMarkup(workoutId int) tgbotapi.InlineKeyboardMarkup {
	callbackData := "leave_" + strconv.Itoa(workoutId)
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{m.Button.Participate("✋Пропущу", callbackData)},
	)
}
