package interfaces

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"telegramStravaBot/domain"
	"time"
)

type YandexWeather interface {
	Init()
	GetForecastRequest() *domain.Weather
	GetForecastText() string
}

type yandexWeather struct {
	bot *tgbotapi.BotAPI
}

func NewYandexWeather(bot *tgbotapi.BotAPI) YandexWeather {
	return &yandexWeather{bot: bot}
}

func (y *yandexWeather) Init() {
	//fmt.Println("test start crontab")
	c := cron.New()
	c.AddFunc("30 5 * * 1,2,4,6", func() {
		fmt.Println("test crontab")
		config := tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: 0}}
		chat, err := y.bot.GetChat(config)
		if err != nil {
			return
		}

		newMessage := tgbotapi.NewMessage(chat.ID, " С добрым утром!\n"+y.GetForecastText())
		y.bot.Send(newMessage)
	})
	c.Start()

}

func (y *yandexWeather) GetForecastRequest() *domain.Weather {
	req, err := http.NewRequest("GET", os.Getenv("YANDEX_URL"), nil)
	if err != nil {
		log.Panic(err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	req.Header.Add("X-Yandex-API-Key", os.Getenv("YANDEX_API_KEY"))
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

	jsonMap := new(domain.Weather)
	jsonErr := json.Unmarshal(body, &jsonMap)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//unixTimeUTC:=time.Unix(jsonMap.Now, 0).Format(time.RFC822)
	//
	//fmt.Printf("Дата: %s \nПогода: %d°, ощущается как: %d°, осадки: %s, ск. ветра: %.1f м/с, влажность: %d",
	//	unixTimeUTC, jsonMap.Fact.Temp, jsonMap.Fact.Feel, jsonMap.Fact.Condition, jsonMap.Fact.WindSpeed, jsonMap.Fact.Humidity)
	return jsonMap
}

func (y *yandexWeather) GetForecastText() string {
	jsonMap := y.GetForecastRequest()
	unixTimeUTC := time.Unix(jsonMap.Now, 0).Format(time.RFC822)

	return fmt.Sprintf("Дата: %s \nПогода: %d°, ощущается как: %d°, осадки: %s, ск. ветра: %.1f м/с, влажность: %d",
		unixTimeUTC, jsonMap.Fact.Temp, jsonMap.Fact.Feel, jsonMap.Fact.Condition, jsonMap.Fact.WindSpeed, jsonMap.Fact.Humidity)
}
