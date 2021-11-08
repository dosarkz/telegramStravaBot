package interfaces

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"telegramStravaBot/domain"
	"time"
)

type TelegramUI interface {
	MainMenu() tgbotapi.ReplyKeyboardMarkup
	StravaInlineButton() tgbotapi.InlineKeyboardButton
	InstaInlineButton() tgbotapi.InlineKeyboardButton
	MetronomeInlineButton() tgbotapi.InlineKeyboardButton
	BotanInlineButton() tgbotapi.InlineKeyboardButton
	MetroInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup
	MarathonInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup
}
type telegramUI struct{}

func (t *telegramUI) MainMenu() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Клуб Любителей Бега MaratHON"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Рейтинг Метронома"),
			tgbotapi.NewKeyboardButton("Расписание"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("О нас"),
			tgbotapi.NewKeyboardButton("Разминка Амосова"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Погода"),
			tgbotapi.NewKeyboardButton("Закрыть меню"),
		),
	)
	return keyboard
}
func (t *telegramUI) StravaInlineButton() tgbotapi.InlineKeyboardButton {
	stravaButtonData := os.Getenv("METRONOME_STRAVA_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Подписаться в Strava", URL: &stravaButtonData}
}
func (t *telegramUI) InstaInlineButton() tgbotapi.InlineKeyboardButton {
	instagramButtonData := os.Getenv("METRONOME_INSTA_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Подписаться в Instagram", URL: &instagramButtonData}
}

func (t *telegramUI) MetronomeInlineButton() tgbotapi.InlineKeyboardButton {
	metronomeButtonData := os.Getenv("METRONOME_TELEGRAM_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Сообщество в Триатлон парке (Metronome)", URL: &metronomeButtonData}
}

func (t *telegramUI) BotanInlineButton() tgbotapi.InlineKeyboardButton {
	botanButtonData := os.Getenv("BOTAN_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Сообщество в Ботаническом саду (Ботаны)", URL: &botanButtonData}
}

func (t *telegramUI) CentralInlineButton() tgbotapi.InlineKeyboardButton {
	centralButtonData := os.Getenv("CENTRAL_URL")
	return tgbotapi.InlineKeyboardButton{Text: "Сообщество в Центральном парке (whatsapp группа)", URL: &centralButtonData}
}

func (t *telegramUI) MetroInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{t.StravaInlineButton(),
		t.InstaInlineButton()})
}

func (t *telegramUI) MarathonInlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{t.MetronomeInlineButton()},
		[]tgbotapi.InlineKeyboardButton{t.BotanInlineButton()},
		[]tgbotapi.InlineKeyboardButton{t.CentralInlineButton()})
}

type TelegramUIRepository struct {
	UI TelegramUI
	YA YandexWeather
}

func NewTelegramUI() TelegramUI {
	return &telegramUI{}
}

func (r TelegramUIRepository) Init(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	for update := range updates {

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "/start":
			msg = getStartMessage(update)
			break
		case "/open":
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
		case "Расписание":
			msg.Text = "Тренировки проводятся на улице в 6:00 утра:\n" +
				"- вторник,\n" +
				"- четверг,\n" +
				"- суббота\n"
			break
		case "Клуб Любителей Бега MaratHON":
			msg.Text = getMarathonInfo()
			msg.ReplyMarkup = r.UI.MarathonInlineKeyboardMarkup()
			break
		case "/club", "О нас":
			msg.Text = getClubInfo()
			msg.ReplyMarkup = r.UI.MetroInlineKeyboardMarkup()
			break
		case "Погода":
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
		case "/hello":
			msg.Text = "Привет, бегун!"
			break
		}

		if msg.Text != update.Message.Text {
			fmt.Printf(msg.Text)
			sendMsg(msg, update, bot)
		}
	}
}

func sendMsg(msg tgbotapi.MessageConfig, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
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
		for i := 0; i < 15; i++ {
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
			"/rating ➡️ Рейтинг Метронома\n"+
			"/club ➡️ Информация о клубе\n")
	msg.ParseMode = "Markdown"
	return msg
}

func getClubInfo() string {
	return "*Рады приветствовать в Metronome team!* \n" +
		"74' выпуск от школы бега Марата Жыланбаева, которая бегает во все времена года!\n" +
		"Друзья, добро пожаловать! Подписывайтесь на наши странички и участвуйте в еженедельном рейтинге бега.\n"
}

func getMarathonInfo() string {
	return "*Marat#ON Клуб Марафонцев* \n" +
		"Клуб Любителей Бега Marat#ON (Марафон) вновь создан в г. Астана в начале 2017 года и объединяет выпускников школы бега Марата Жыланбаева.\n Мастер спорта международного класса, ультрамарафонец, первый и единственный атлет в истории человечества, в одиночку пробежавший крупнейшие пустыни Азии, Африки, Австралии и Америки.\n Установил несколько мировых рекордов, семь из них занесены в Книгу рекордов Гиннеса.\n Большая часть мировых рекордов, установленных Жыланбаевым в начале 1990-х годов остаются по-прежнему не превзойденными.\n"
}
