package interfaces

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"telegramStravaBot/domain"
	user "telegramStravaBot/domain/users"
	"telegramStravaBot/domain/workouts"
	"time"
)

type TelegramUI interface {
	MainMenu() tgbotapi.ReplyKeyboardMarkup
	StravaInlineButton() tgbotapi.InlineKeyboardButton
	InstaInlineButton() tgbotapi.InlineKeyboardButton
	Participate(text string, callback string) tgbotapi.InlineKeyboardButton
	MetronomeInlineButton() tgbotapi.InlineKeyboardButton
	BotanInlineButton() tgbotapi.InlineKeyboardButton
	MetroInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup
	MarathonInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup
	AppointmentKeyboardMarkup(workoutId int) tgbotapi.InlineKeyboardMarkup
	AppointmentDoneKeyboardMarkup(workoutId int) tgbotapi.InlineKeyboardMarkup
	HideMenu() tgbotapi.ReplyKeyboardRemove
	MarathonText() string
}
type telegramUI struct{}

func (t *telegramUI) MainMenu() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Клуб Любителей Бега MaratHON"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Рейтинг Метронома"),
			tgbotapi.NewKeyboardButton("Запись на тренировку"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Разминка Амосова"),
			tgbotapi.NewKeyboardButton("Погода"),
		),
	)
	return keyboard
}

func (t *telegramUI) MarathonText() string {
	return "*Marat#ON Клуб Марафонцев* \n" +
		"Клуб Любителей Бега Marat#ON (Марафон) вновь создан в г. Астана в начале 2017 года" +
		" и объединяет выпускников школы бега Марата Жыланбаева.\n Мастер спорта международного класса," +
		" ультрамарафонец, первый и единственный атлет в истории человечества, в одиночку пробежавший крупнейшие" +
		" пустыни Азии, Африки, Австралии и Америки.\n Установил несколько мировых рекордов, семь из них занесены в" +
		" Книгу рекордов Гиннеса.\n Большая часть мировых рекордов, установленных Жыланбаевым в начале 1990-х годов" +
		" остаются по-прежнему не превзойденными.\n" +
		"\n"
}

func (t *telegramUI) HideMenu() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
}

func (t *telegramUI) StravaInlineButton() tgbotapi.InlineKeyboardButton {
	stravaButtonData := os.Getenv("METRONOME_STRAVA_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Подписаться в Strava", URL: &stravaButtonData}
}

func (t *telegramUI) AlmatyInlineButton() tgbotapi.InlineKeyboardButton {
	almatyButtonData := os.Getenv("ALMATY_GROUP_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Алматинское сообщество бега", URL: &almatyButtonData}
}

func (t *telegramUI) InstaInlineButton() tgbotapi.InlineKeyboardButton {
	instagramButtonData := os.Getenv("METRONOME_INSTA_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Подписаться в Instagram", URL: &instagramButtonData}
}

func (t *telegramUI) MetronomeInlineButton() tgbotapi.InlineKeyboardButton {
	metronomeButtonData := os.Getenv("METRONOME_TELEGRAM_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Сообщество бега в Триатлон парке (Metronome)", URL: &metronomeButtonData}
}

func (t *telegramUI) Participate(text string, callback string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{Text: text, CallbackData: &callback}
}

func (t *telegramUI) BotanInlineButton() tgbotapi.InlineKeyboardButton {
	botanButtonData := os.Getenv("BOTAN_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Сообщество бега в Ботаническом саду (Ботаны)", URL: &botanButtonData}
}

func (t *telegramUI) CentralInlineButton() tgbotapi.InlineKeyboardButton {
	centralButtonData := os.Getenv("CENTRAL_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Сообщество бега в Центральном парке (whatsapp группа)", URL: &centralButtonData}
}

func (t *telegramUI) MetroInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{t.StravaInlineButton(),
		t.InstaInlineButton()})
}
func (t *telegramUI) AlmatyRunningGroups() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{t.AlmatyInlineButton()})
}

func (t *telegramUI) MarathonInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{t.MetronomeInlineButton()},
		[]tgbotapi.InlineKeyboardButton{t.BotanInlineButton()},
		[]tgbotapi.InlineKeyboardButton{t.CentralInlineButton()},
		[]tgbotapi.InlineKeyboardButton{t.AlmatyInlineButton()},
	)
}

func (t *telegramUI) AppointmentKeyboardMarkup(workoutId int) tgbotapi.InlineKeyboardMarkup {
	callbackData := "appointment_" + strconv.Itoa(workoutId)
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{t.Participate("✅Приму участие", callbackData)},
	)
}

func (t *telegramUI) AppointmentDoneKeyboardMarkup(workoutId int) tgbotapi.InlineKeyboardMarkup {
	callbackData := "leave_" + strconv.Itoa(workoutId)
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{t.Participate("✋Пропущу", callbackData)},
	)
}

type TelegramUIRepository struct {
	UI      TelegramUI
	YA      YandexWeather
	User    user.UserService
	Workout workouts.WorkoutService
	Bot     *tgbotapi.BotAPI
}

func NewTelegramUI() TelegramUI {
	return &telegramUI{}
}

func (r TelegramUIRepository) Init() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := r.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := r.Bot.Request(callback); err != nil {
				panic(err)
			}
			callbackData := strings.Split(update.CallbackQuery.Data, "_")
			//fmt.Println(callbackData[0])
			workoutId, _ := strconv.Atoi(callbackData[1])

			switch callbackData[0] {
			case "appointment":
				msgText := getAppointmentText(update, 1, r, workoutId)
				answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					msgText, r.UI.AppointmentDoneKeyboardMarkup(workoutId))
				answer.ParseMode = "markdown"
				_, err := r.Bot.Send(answer)
				if err != nil {
					return
				}
				break
			case "leave":
				msgText := getAppointmentText(update, 0, r, workoutId)
				answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID,
					msgText, r.UI.AppointmentKeyboardMarkup(workoutId))
				answer.ParseMode = "markdown"
				_, err := r.Bot.Send(answer)
				if err != nil {
					return
				}
				break
			}
		}
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "/start":
			msg = getStartMessage(update)
			break
		case "/open", "menu":
			msg.ReplyMarkup = r.UI.MainMenu()
			msg.Text = "Открыто меню"
			break
		case "/close", "Закрыть меню":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			msg.Text = "Закрыто"
			break
		case "/rating", "Рейтинг Метронома":
			msg.Text = getRatingClub()
			break
		case "/run", "Запись на тренировку":
			appointmentToRunning(r, update)
			break
		case "Клуб Любителей Бега MaratHON":
			msg.Text = r.UI.MarathonText()
			msg.ReplyMarkup = r.UI.MarathonInlineKeyboardMarkup()
			break
		case "Свернуть меню":
			msg.Text = "ok"
			msg.ReplyMarkup = r.UI.HideMenu()
			break
		case "Погода", "/weather":
			msg.Text = r.YA.GetForecastText()
			break
		case "Разминка Амосова":
			msg.Text = "Уникальная программа разминки Амосова от Марата Толегеновича, делайте ее ежедневно и будете здоровы!\n" +
				"* Каждое упражнение выполняется по 100 раз!* \n" +
				"- Руки на пояс и движение корпусом влево и вправо \n" +
				"- Движение корпусом вниз и вверх \n" +
				"- Сгибание и разгибание рук к центру \n" +
				"- Ноги вместе, сгибание и полуприсед вперед и назад до начала стопы \n" +
				"- Сгибание и полуприсед влево и вправо \n" +
				"- Круговые движения ног по часовой стрелке и против. \n" +
				"- Выпады. \n"
			break
		case "Назад":
			msg.ReplyMarkup = r.UI.MainMenu()
			msg.Text = "Главное меню"
			break
		case "/hello":
			msg.Text = "Привет, бегун!"
			break
		}

		if msg.Text != update.Message.Text {
			fmt.Printf(msg.Text)
			replyMessage(msg, update, r.Bot)
		}
	}
}

func getAppointmentText(update tgbotapi.Update, typeId int, r TelegramUIRepository, workoutId int) string {
	getUser, err := r.User.FindUserByTelegramId(update.CallbackQuery.From.ID)
	if err != nil {
		newUser := user.User{
			Username:   update.CallbackQuery.From.FirstName + " " + update.CallbackQuery.From.LastName,
			TelegramId: update.CallbackQuery.From.ID,
		}
		getUser, err = r.User.CreateUser(&newUser)
		if err != nil {
			log.Panic(err)
		}
	}

	switch typeId {
	case 1:
		registerUserWorkout(getUser, workoutId, r.Workout)
		break
	case 0:
		leaveUserWorkout(getUser, workoutId, r.Workout)
		break
	}
	workout, err := r.Workout.ReadWorkout(workoutId)
	if err != nil {
		log.Panic(err)
	}

	text := fmt.Sprintf("🔥Тренировка № %d\n 🏃‍♀ 🏃 %s\n %s\n %s\n %s", workoutId, workout.CreatedAt.Format(time.RFC822), workout.Title, workout.Description,
		getWorkoutUserList(workoutId, r))
	return text
}

func getWorkoutUserList(workoutId int, r TelegramUIRepository) string {
	wc, err := r.Workout.ListWorkoutMembers(workoutId)
	if err != nil {
		log.Panic(err)
	}

	if len(wc) > 0 {
		msg := "** Участники **\n"
		for i := 0; i < len(wc); i++ {
			readUser, err := r.User.ReadUser(wc[i].UserID)
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

func registerUserWorkout(user *user.User, workoutId int, repository workouts.WorkoutRepository) *workouts.WorkoutUser {
	getWorkoutUser, err := repository.FindBy(user.Id, workoutId)

	if err != nil {
		newWorkoutUser := workouts.WorkoutUser{
			UserID:    user.Id,
			WorkoutId: uint(workoutId),
		}
		getWorkoutUser, err = repository.CreateWorkoutUser(&newWorkoutUser)
		if err != nil {
			log.Panic(err)
		}
	}

	return getWorkoutUser
}

func leaveUserWorkout(user *user.User, workoutId int, repository workouts.WorkoutRepository) {
	getWorkoutUser, err := repository.FindBy(user.Id, workoutId)
	if err != nil {
		fmt.Println(err)
		//log.Panic(err)
		return
	}
	fmt.Println(getWorkoutUser)
	_, err = repository.Delete(getWorkoutUser)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func replyMessage(msg tgbotapi.MessageConfig, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

func appointmentToRunning(r TelegramUIRepository, update tgbotapi.Update) {
	ws, err := r.Workout.ListWorkouts()
	if err != nil {
		log.Panic(err)
	}

	var responseItems = make([]workouts.WorkoutResponse, len(ws))

	for i, element := range ws {
		responseItems[i] = *workouts.ToResponseModel(&element)
	}

	if len(responseItems) == 0 {
		_, err := r.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "🎽Ближайших тренировок не наблюдается, "+
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
		newMessage.ReplyMarkup = r.UI.AppointmentKeyboardMarkup(responseItems[i].Id)
		newMessage.ParseMode = "markdown"
		_, err := r.Bot.Send(newMessage)
		if err != nil {
			return
		}
	}
}

func getRatingClub() string {
	currentTime := time.Now().Format(time.RFC822)
	message := "🏆Рейтинг Метронома 🏃‍♀️🏃 на этой неделе от " + currentTime + "\n"
	message += "\n"
	req, err := http.NewRequest("GET", "https://www.strava.com/clubs/540448/leaderboard", nil)
	if err != nil {
		log.Panic(err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	req.Header.Add("x-requested-with", `XMLHttpRequest`)
	req.Header.Add("Accept", `text/javascript, application/javascript, application/ecmascript, application/x-ecmascript`)

	resp, err := client.Do(req)

	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Panic(err)
			}
		}(resp.Body)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonMap := make(map[string][]domain.Rating)
	jsonErr := json.Unmarshal(body, &jsonMap)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	for _, items := range jsonMap {
		for i := 0; i < len(items); i++ {
			if i >= 15 {
				break
			}
			athleteLink := fmt.Sprintf("https://www.strava.com/athletes/%d", items[i].AthleteId)
			message += fmt.Sprintf("%d) [%s %s](%s) - расстояние: %.2f км, забеги: %d, самый длинный: %.2f км, ср.темп: %.2f /км \n", items[i].Rank, items[i].AthleteFirstname, items[i].AthleteLastname, athleteLink,
				items[i].Distance/1000, items[i].NumActivities, items[i].BestActivitiesDistance/1000, getPace(items[i].MovingTime, items[i].Distance))
		}
	}
	message += "\n"
	message += "**Хотите участвовать в рейтинге ‍🚀?** \n Подписывайтесь на страницу в [STRAVA](https://www.strava.com/clubs/540448) и вы автоматически будете в нашем списке 😀👍"

	return message
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

func getStartMessage(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"*Вас приветствует бот Metronome* 😃🖐"+
			"\n"+
			"/open ➡️ Открыть главное меню\n"+
			"/close ➡️ Закрыть меню\n"+
			"/rating ➡️ Рейтинг Метронома\n")
	msg.ParseMode = "Markdown"
	return msg
}
