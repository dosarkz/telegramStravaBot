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

	if update.CallbackQuery.Data == "workout_complete" {
		errs := s.Redis.Set("makeWorkout", 0, 0).Err()
		if errs != nil {
			panic(errs)
		}
		return
	}
	callbackData := strings.Split(update.CallbackQuery.Data, "_")
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
