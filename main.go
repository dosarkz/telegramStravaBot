package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
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
	// 	tgbotapi.NewKeyboardButtonRow(
	// 		tgbotapi.NewKeyboardButton("Местоположение точки сбора"),
	// 		tgbotapi.NewKeyboardButton("СБУ И ОФП"),
	// 	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Закрыть меню"),
	),
)

var instagramButtonData = "https://www.instagram.com/metronom_team"
var stravaButtonData = "https://www.strava.com/clubs/540448"
var metronomeButtonData = "https://t.me/joinchat/VoeA783qZuIBa4um"
var botanButtonData = "https://t.me/botandostar"
var centralButtonData = "https://chat.whatsapp.com/LdfZSnyInE7F7PrAfpqaXj"
var stravaButton = tgbotapi.InlineKeyboardButton{Text: "Подписаться в Strava", URL: &stravaButtonData}
var instaButton = tgbotapi.InlineKeyboardButton{Text: "Подписаться в Instagram", URL: &instagramButtonData}

var metronomeButton = tgbotapi.InlineKeyboardButton{Text: "Сообщество в Триатлон парке (Metronome)", URL: &metronomeButtonData}
var botanButton = tgbotapi.InlineKeyboardButton{Text: "Сообщество в Ботаническом саду (Ботаны)", URL: &botanButtonData}
var centralButton = tgbotapi.InlineKeyboardButton{Text: "Сообщество в Центральном парке (whatsapp группа)", URL: &centralButtonData}
var metroKeyBoard = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{stravaButton, instaButton})
var marathonKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{metronomeButton},
	[]tgbotapi.InlineKeyboardButton{botanButton},
	[]tgbotapi.InlineKeyboardButton{centralButton})

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	c := cron.New()
	c.AddFunc("30 5 * * 2,4,6", func() {
		fmt.Println("Every hour on the half hour")
		config := tgbotapi.ChatConfig{ChatID: 222288800}
		chat, err := bot.GetChat(config)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("chat id is:  %s", chat)
		newMessage := tgbotapi.NewMessage(chat.ID, "Qairly Tan, Dostar!  ")
		bot.Send(newMessage)
	})
	c.Start()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "/start":
			msg = getStartMessage(update)
			break
		case "/open":
			msg.ReplyMarkup = numericKeyboard
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
			msg.ReplyMarkup = marathonKeyBoard
			break
		case "/club", "О нас":
			msg.Text = getClubInfo()
			msg.ReplyMarkup = metroKeyBoard
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

	jsonMap := make(map[string][]Rating)
	jsonErr := json.Unmarshal(body, &jsonMap)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	for _, items := range jsonMap {
		for i := 0; i < len(items); i++ {
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

type Rating struct {
	Distance      float32 `json:"distance"`
	NumActivities int     `json:"num_activities"`
	//elevGain   float32
	MovingTime             int     `json:"moving_time"`
	Velocity               float32 `json:"velocity"`
	BestActivitiesDistance float32 `json:"best_activities_distance"`
	Rank                   int     `json:"rank"`
	AthleteFirstname       string  `json:"athlete_firstname"`
	AthleteId              int     `json:"athlete_id"`
	AthleteLastname        string  `json:"athlete_lastname"`
	//athletePictureUrl string
	//athleteMemberType string
}
