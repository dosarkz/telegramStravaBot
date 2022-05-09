package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

func (a UIActionService) callbackQuery(update tgbotapi.Update, s *UIService) {
	if update.CallbackQuery == nil {
		return
	}
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := a.Bot.Request(callback); err != nil {
		log.Panic(err)
		return
	}

	switch update.CallbackQuery.Data {
	case "workout_complete":
		errs := s.Redis.Set("makeWorkout", 0, 0).Err()
		if errs != nil {
			log.Panic(errs)
			return
		}
		break
	case "update_metro":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
		waitingMsg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			update.CallbackQuery.Message.Text, s.Menu.HeroUpdatingButtonKeyboard())
		s.Action.Bot.Send(waitingMsg)

		metroMessage := getRatingMessage(msg)
		upd := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			metroMessage.Text, s.Menu.MetroUpdateButtonKeyboard())
		upd.ParseMode = "markdown"
		s.Action.Bot.Send(upd)
		break
	case "update_hero":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
		waitingMsg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			update.CallbackQuery.Message.Text, s.Menu.HeroUpdatingButtonKeyboard())
		s.Action.Bot.Send(waitingMsg)

		heroMessage := getHeroByDay(msg)
		upd := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			heroMessage.Text, s.Menu.HeroUpdateButtonKeyboard())
		upd.ParseMode = "markdown"
		s.Action.Bot.Send(upd)
		break
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
