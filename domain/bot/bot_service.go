package bot

import (
	"encoding/json"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"telegramStravaBot/domain"
	"telegramStravaBot/domain/workouts"
)

type UIService struct {
	Menu   *UIMenuService
	Action *UIActionService
	Repos  *domain.Repositories
	Redis  *redis.Client
}

func NewUIService(service UIActionService, repos *domain.Repositories, redis *redis.Client) *UIService {
	return &UIService{
		Action: &service,
		Menu:   &UIMenuService{Button: &UIButtonService{}},
		Repos:  repos,
		Redis:  redis,
	}
}

func (s UIService) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := s.Action.Bot.GetUpdatesChan(u)

	go s.Repos.UserRepository.FindWorkoutsAndSaveScore()

	for update := range updates {
		s.Action.callbackQuery(update, &s)
		mwk := []byte(s.Redis.Get("makeWorkout").Val())

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		isGroup := checkIsGroup(update, msg, s.Action.Bot)

		switch update.Message.Text {
		case "⚡ Рейтинг метронома":
			msg = getRatingMessage(msg)
			break
		case "✅ Записаться":
			appointmentToRunning(&s, update)
			break
		case "💥 Герой дня":
			msg = getHeroByDay(msg)
			break
		case "🏃 Клуб Любителей Бега MaratHON":
			msg = getClubMessage(msg, s.Menu)
			break
		case "😊 Разминка Амосова":
			msg.Text = amosovMessageText()
			break
		case "☂ Погода":
			msg.Text = s.Action.YA.GetForecastText()
			break
		case "➕Добавить тренировку":
			if isGroup {
				continue
			}
			msg = newTraining(msg, update, s.Redis)
			msg.ReplyMarkup = s.Menu.CreateWorkoutKeyboard()
			break
		}

		switch update.Message.Command() {
		case "start":
			msg = getStartMessage(update)
			break
		case "open":
			msg = getOpenMessage(msg, s.Menu)
			break
		case "close":
			msg = getCloseMessage(msg)
			break
		case "rating":
			msg = getRatingMessage(msg)
			break
		case "newWorkout":
			if isGroup {
				continue
			}
			msg = newTraining(msg, update, s.Redis)
			msg.ReplyMarkup = s.Menu.CreateWorkoutKeyboard()
			break
		case "skipNewWorkout":
			if isGroup {
				continue
			}
			msg.Text = "Создание новой тренировки прервано успешно."
			err := s.Redis.Set("makeWorkout", 0, 0).Err()
			if err != nil {
				log.Panic(err)
			}
			break
		case "deleteNewWorkout":
			if isGroup {
				continue
			}
			msg.Text = "Отправьте id тренировки следующим сообщением для удаления записи."
			bJson, err := json.Marshal(&workouts.WorkoutStatus{UserId: update.Message.From.ID,
				DeleteStatus: 1})
			if err != nil {
				log.Panic(err)
			}
			err = s.Redis.Set("makeWorkout", bJson, 0).Err()
			if err != nil {
				log.Panic(err)
			}
			break
		case "run":
			appointmentToRunning(&s, update)
			break
		case "club":
			msg = getClubMessage(msg, s.Menu)
			break
		case "weather":
			msg.Text = s.Action.YA.GetForecastText()
			break
		case "amosov":
			msg.Text = amosovMessageText()
			break
		case "back":
			msg.ReplyMarkup = s.Menu.MainMenu()
			msg.Text = "Главное меню"
			break
		}

		data := workouts.WorkoutStatus{}
		json.Unmarshal(mwk, &data)

		if !isGroup && data.UserId == update.Message.From.ID {
			if data.CreateStatus != 0 {
				wErr, workout := s.Repos.WorkoutRepository.CallbackNewWorkout(update)

				if wErr {
					msg.Text = "Шаблон тренировки введен некорректно. Пожалуйста попробуйте изменить текст и выполнить команду занова."
				} else {
					msg.Text = "Успешно. Тренировка сохранена под №" + strconv.Itoa(workout.Id)
					err := s.Redis.Set("makeWorkout", 0, 0).Err()
					if err != nil {
						log.Panic(err)
					}
				}
			}

			if data.DeleteStatus != 0 {
				wErr := s.Repos.WorkoutRepository.CallbackDeleteWorkout(update)
				if wErr {
					s.Redis.Set("makeWorkout", 0, 0)
					msg.Text = "Тренировка удалена"
				}
			}
		}

		if msg.Text != update.Message.Text {
			replyMessage(msg, update, s.Action.Bot)
		}
	}
}

func newTraining(msg tgbotapi.MessageConfig, update tgbotapi.Update, redis *redis.Client) tgbotapi.MessageConfig {
	msg.Text = getWorkoutNewMessage()
	bJson, err := json.Marshal(&workouts.WorkoutStatus{UserId: update.Message.From.ID,
		CreateStatus: 1})

	err = redis.Set("makeWorkout", bJson, 0).Err()
	if err != nil {
		log.Panic(err)
	}

	return msg
}

func checkIsGroup(update tgbotapi.Update, msg tgbotapi.MessageConfig, Bot *tgbotapi.BotAPI) bool {
	if update.Message.Chat.Type == "group" {
		msg.Text = "Данная функция доступна только в персональном чате с ботом!"
		replyMessage(msg, update, Bot)
		return true
	}
	return false
}
