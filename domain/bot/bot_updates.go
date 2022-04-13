package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a UIActionService) callbackQuery(update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := a.Bot.Request(callback); err != nil {
		panic(err)
	}

	fmt.Println(update.CallbackQuery.Data)

	//switch callbackData[0] {
	//case "appointment":
	//	joinWorkout()
	//	break
	//case "leave":
	//	leaveWorkout()
	//	break
	//}
}
