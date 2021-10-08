package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
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
		tgbotapi.NewKeyboardButton("Местоположение точки сбора"),
		tgbotapi.NewKeyboardButton("СБУ И ОФП"),
	),
	tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("Закрыть меню"),
    ),
)

var stravaButtonData = "👍🏻"
var metroKeyBoard  = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{tgbotapi.InlineKeyboardButton{Text: "Страница в Strava", CallbackData: &stravaButtonData}})

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
            case "/club":
                msg.Text = getClubInfo()
                msg.ReplyMarkup = metroKeyBoard
                break
            case "/hello":
                msg.Text = "Привет, бегун!"
				break
            default:
				msg.Text = "Ой, кажется что-то пошло не так."
		}

		sendMsg(msg, update, bot)
	}
}

func sendMsg(msg tgbotapi.MessageConfig, update tgbotapi.Update, bot *tgbotapi.BotAPI )  {
	msg.ReplyToMessageID = update.Message.MessageID
    msg.ParseMode = "markdown"
    bot.Send(msg)
}

func getRatingClub() string {
	currentTime := time.Now().Format("01-02-2006")
	message := "Рейтинг Метронома на этой неделе от "+currentTime+"\n"
	req, err := http.NewRequest("GET", "https://www.strava.com/clubs/540448/leaderboard", nil)
	if err != nil {
		log.Panic(err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	req.Header.Add("x-requested-with", `XMLHttpRequest`)
	req.Header.Add("Accept", `text/javascript, application/javascript, application/ecmascript, application/x-ecmascript`)

	resp, err := client.Do(req)

	if resp.Body != nil {
		defer resp.Body.Close()
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
		for i:=0; i<len(items); i++ {
			athleteLink := fmt.Sprintf("https://www.strava.com/athletes/%d", items[i].AthleteId)
			message += fmt.Sprintf("%d) [%s %s](%s) - расстояние: %1.f км, забеги: %d, самый длинный: %1.f км \n", items[i].Rank, items[i].AthleteFirstname, items[i].AthleteLastname, athleteLink,
				items[i].Distance / 1000, items[i].NumActivities, items[i].BestActivitiesDistance / 1000)
		}
	}

	return message
}

func getStartMessage(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"*Вас приветствует бот Metronome* 😃🖐"+
			"\n"+
			"/open ➡️ Открыть главное меню\n"+
			"/close ➡️ Закрыть меню\n"+
			"/rating ➡️ Рейтинг Метронома\n" +
			"/club ➡️ Информация о клубе\n")
	msg.ParseMode = "Markdown"
	return msg
}

func getClubInfo() string{
    return "*Рады приветствовать в Metronome team!* \n" +
    	"74' выпуск от школы бега Марата Жыланбаева, которая бегает во все времена года!\n" +
    	"Друзья, добро пожаловать!\n"
}

type Rates struct {
	data []Rating
}

type Rating struct {
	Distance float32 `json:"distance"`
	NumActivities int `json:"num_activities"`
	//elevGain   float32
	MovingTime int `json:"moving_time"`
	Velocity               float32 `json:"velocity"`
	BestActivitiesDistance float32 `json:"best_activities_distance"`
	Rank             int `json:"rank"`
	AthleteFirstname string `json:"athlete_firstname"`
	AthleteId int `json:"athlete_id"`
	AthleteLastname     string `json:"athlete_lastname"`
	//athletePictureUrl string
	//athleteMemberType string
}