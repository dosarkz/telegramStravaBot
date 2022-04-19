package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

func (a UIActionService) callbackQuery(update tgbotapi.Update, s *UIService) {
	if update.CallbackQuery == nil {
		return
	}
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := a.Bot.Request(callback); err != nil {
		panic(err)
	}
	callbackData := strings.Split(update.CallbackQuery.Data, "_")
	//fmt.Println(callbackData[0])
	workoutId, _ := strconv.Atoi(callbackData[1])

	switch callbackData[0] {
	case "appointment":
		join(workoutId, update, s)
		break
	case "leave":
		leave(workoutId, update, s)
		break
	}
}
