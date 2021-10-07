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

func main() {
	err := godotenv.Load()
	message := ""
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

		switch update.Message.Text {
		case "/rating":
			message = getRatingClub()
			break
		case "/hello":
			message = "Привет, бегун!"
		default:
			message = "Хм... бот таких команд не знает.\n Посмотреть рейтинг по клубу: /rating"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "markdown"
		bot.Send(msg)
	}
}

func getRatingClub() string {
	message := "Рейтинг Метронома на этой неделе\n"
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