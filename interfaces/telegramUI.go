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
	"telegramStravaBot/domain"
	user "telegramStravaBot/domain/users"
	"telegramStravaBot/domain/workouts"
	"time"
)

type TelegramUI interface {
	MainMenu() tgbotapi.ReplyKeyboardMarkup
	StravaInlineButton() tgbotapi.InlineKeyboardButton
	InstaInlineButton() tgbotapi.InlineKeyboardButton
	Participate(text string, callback *string) tgbotapi.InlineKeyboardButton
	MetronomeInlineButton() tgbotapi.InlineKeyboardButton
	BotanInlineButton() tgbotapi.InlineKeyboardButton
	MetroInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup
	MarathonInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup
	AppointmentKeyboardMarkup() tgbotapi.InlineKeyboardMarkup
	AppointmentDoneKeyboardMarkup() tgbotapi.InlineKeyboardMarkup
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

func (t *telegramUI) Participate(text string, callback *string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{Text: text, CallbackData: callback}
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

func (t *telegramUI) AppointmentKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {
	callback := "appointment"
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{t.Participate("Принять участие", &callback)},
	)
}

func (t *telegramUI) AppointmentDoneKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {
	callback := "do_not_participate"
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{t.Participate("Больше не участвовать", &callback)},
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

			switch update.CallbackQuery.Data {
			case "appointment":
				msgText := update.CallbackQuery.Message.ReplyToMessage.From.UserName
				answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID,
					msgText, r.UI.AppointmentDoneKeyboardMarkup())
				r.Bot.Send(answer)
				break
			case "do_not_participate":
				msgText := update.CallbackQuery.Message.ReplyToMessage.From.UserName
				answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID,
					msgText, r.UI.AppointmentKeyboardMarkup())
				r.Bot.Send(answer)
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
		case "Запись на тренировку":
			appointmentToRunning(r, update.Message.Chat.ID)
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

func replyMessage(msg tgbotapi.MessageConfig, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

func appointmentToRunning(r TelegramUIRepository, chatId int64) {
	ws, err := r.Workout.ListWorkouts()
	if err != nil {
		log.Panic(err)
	}

	var responseItems = make([]workouts.WorkoutResponse, len(ws))

	for i, element := range ws {
		responseItems[i] = *workouts.ToResponseModel(&element)
	}

	for i := 0; i < len(responseItems); i++ {
		msg := fmt.Sprintf("%s\n %s\n Дата: %s\n", responseItems[i].Title, responseItems[i].Description,
			responseItems[i].CreatedAt)

		newMessage := tgbotapi.NewMessage(chatId, msg)
		newMessage.ReplyMarkup = r.UI.AppointmentKeyboardMarkup()
		newMessage.ParseMode = "markdown"
		r.Bot.Send(newMessage)
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
