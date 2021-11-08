package interfaces

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron"
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
	c := cron.New()
	c.AddFunc("30 5 * * 2,4,6", func() {
		config := tgbotapi.ChatConfig{ChatID: -1001451720943}
		chat, err := y.bot.GetChat(config)
		if err != nil {
			log.Panic(err)
		}

		newMessage := tgbotapi.NewMessage(chat.ID, " С добрым утречком тебя,\n"+
			"Улыбнись скорее,\n"+
			"Легкого желаю дня☀,\n"+
			"Быть тебе бодрее!\n"+
			"\n"+
			"Всюду и везде успеть,\n"+
			"Чаще улыбаться,\n"+
			"А еще не уставать,\n"+
			"Жизнью наслаждаться🤗!\n\n"+y.GetForecastText())
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
