package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegramStravaBot/domain"
)

type UIService struct {
	Menu   *UIMenuService
	Action *UIActionService
	Repos  *domain.Repositories
}

func NewUIService(service UIActionService, repos *domain.Repositories) *UIService {
	return &UIService{
		Action: &service,
		Menu:   &UIMenuService{Button: &UIButtonService{}},
		Repos:  repos,
	}
}

func (s UIService) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := s.Action.Bot.GetUpdatesChan(u)

	for update := range updates {
		s.Action.callbackQuery(update)

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "Рейтинг Метронома":
			msg = getRatingMessage(msg)
			break
		case "Запись на тренировку":
			appointmentToRunning(&s, update)
			break
		case "Клуб Любителей Бега MaratHON":
			msg = getClubMessage(msg, s.Menu)
			break
		case "Разминка Амосова":
			msg.Text = amosovMessageText()
			break
		case "Погода":
			msg.Text = s.Action.YA.GetForecastText()
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
		case "new_workout":
			if update.Message.Chat.Type == "group" {
				msg.Text = "Добавлять тренировки можно только персональном чате с ботом!"
				break
			}
			msg.Text = getWorkoutNewMessage()
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

		//if isCreateWorkout != 0 {
		//	fmt.Printf("creation a workout\n")
		//	newText := strings.Split(update.Message.Text, "\n")
		//	if len(newText) < 2 {
		//		msg.Text = "Шаблон тренировки введен некорректно. Пожалуйста попробуйте изменить текст и выполнить команду занова."
		//		replyMessage(msg, update, r.Bot)
		//		continue
		//	}
		//
		//	date, err := time.Parse("2.1.2006 15:04", newText[2])
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//
		//	fmt.Printf("date %s\n", date)
		//	wk := &workouts.Workout{Title: newText[0], Description: newText[1], CreatedAt: date}
		//	fmt.Println(wk)
		//}

		if msg.Text != update.Message.Text {
			replyMessage(msg, update, s.Action.Bot)
		}
	}
}
