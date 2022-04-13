package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math"
	"telegramStravaBot/domain/workouts"
	"time"
)

type Participation interface {
	join()
	leave()
}

//func join(update tgbotapi.Update) {
//	msgText := getAppointmentText(update)
//	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
//		update.CallbackQuery.Message.MessageID,
//		msgText, a.UIService.Menu.AppointmentDoneKeyboardMarkup(workoutId))
//	answer.ParseMode = "markdown"
//	_, err := a.Bot.Send(answer)
//	if err != nil {
//		return
//	}
//}
//
//func leave(update tgbotapi.Update) {
//	msgText := getAppointmentText(update, 0, a, workoutId)
//	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID,
//		msgText, a.UIService.Menu.AppointmentKeyboardMarkup(workoutId))
//	answer.ParseMode = "markdown"
//	_, err := a.Bot.Send(answer)
//	if err != nil {
//		return
//	}
//}
//
//func getMember() {
//	u, err := r.User.FindUserByTelegramId(update.CallbackQuery.From.ID)
//	if err != nil {
//		newUser := user.User{
//			Username:   update.CallbackQuery.From.FirstName + " " + update.CallbackQuery.From.LastName,
//			TelegramId: update.CallbackQuery.From.ID,
//		}
//		u, err = r.User.CreateUser(&newUser)
//		if err != nil {
//			log.Panic(err)
//		}
//	}
//	return u
//}
//
//func getAppointmentText(update tgbotapi.Update, typeId int, r TelegramUIRepository, workoutId int) string {
//
//	us := getMember()
//
//	switch typeId {
//	case 1:
//		registerUserWorkout(getUser, workoutId, r.Workout)
//		break
//	case 0:
//		leaveUserWorkout(getUser, workoutId, r.Workout)
//		break
//	}
//	workout, err := r.Workout.ReadWorkout(workoutId)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	text := fmt.Sprintf("🔥Тренировка № %d\n 🏃‍♀ 🏃 %s\n %s\n %s\n %s", workoutId, workout.CreatedAt.Format(time.RFC822), workout.Title, workout.Description,
//		getWorkoutUserList(workoutId, r))
//	return text
//}
//
func getWorkoutUserList(workoutId int, r *UIService) string {
	wc, err := r.Repos.ListWorkoutMembers(workoutId)
	if err != nil {
		log.Panic(err)
	}

	if len(wc) > 0 {
		msg := "** Участники **\n"
		for i := 0; i < len(wc); i++ {
			readUser, err := r.Repos.ReadUser(wc[i].UserID)
			if err != nil {
				log.Panic(err)
			}
			msg += fmt.Sprintf("%d. %s\n", i+1, readUser.Username)
		}
		return msg
	} else {
		return ""
	}
}

func appointmentToRunning(r *UIService, update tgbotapi.Update) {
	ws, err := r.Repos.WorkoutRepository.ListWorkouts()
	if err != nil {
		log.Panic(err)
	}

	var responseItems = make([]workouts.WorkoutResponse, len(ws))

	for i, element := range ws {
		responseItems[i] = *workouts.ToResponseModel(&element)
	}

	if len(responseItems) == 0 {
		_, err := r.Action.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "🎽Ближайших тренировок не наблюдается, "+
			"отдыхай и восстанавливайся😴."))
		if err != nil {
			return
		}
	}

	for i := 0; i < len(responseItems); i++ {
		msg := fmt.Sprintf("🔥Тренировка № %d\n 🏃‍♀ 🏃 %s\n %s\n %s\n %s\n", responseItems[i].Id,
			responseItems[i].CreatedAt.Format(time.RFC822),
			responseItems[i].Title, responseItems[i].Description, getWorkoutUserList(responseItems[i].Id, r))

		newMessage := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		newMessage.ReplyMarkup = r.Menu.AppointmentKeyboardMarkup(responseItems[i].Id)
		newMessage.ParseMode = "markdown"
		_, err := r.Action.Bot.Send(newMessage)
		if err != nil {
			return
		}
	}
}

func secondsToMinutes(inSeconds int) float64 {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	return float64(minutes + (seconds / 100))
}

func getPace(movingTime int, distance float32) float64 {
	intPace, float := math.Modf(secondsToMinutes(movingTime) / float64(distance/1000))
	return intPace + (float * 60 / 100)
}
