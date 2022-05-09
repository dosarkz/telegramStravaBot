package bot

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"telegramStravaBot/config"
	"telegramStravaBot/domain/strava/entities"
	"time"
)

func getStartMessage(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, getStartMessageText())
	msg.ParseMode = "Markdown"
	return msg
}

func getOpenMessage(msg tgbotapi.MessageConfig, service *UIMenuService) tgbotapi.MessageConfig {
	msg.ReplyMarkup = service.MainMenu()
	msg.Text = "Открыто меню"
	return msg
}

func getCloseMessage(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.Text = "Закрыто"
	return msg
}

func getClubMessage(msg tgbotapi.MessageConfig, service *UIMenuService) tgbotapi.MessageConfig {
	msg.Text = getClubMessageText()
	msg.ReplyMarkup = service.MarathonInlineKeyboardMarkup()
	return msg
}

func getRatingMessage(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
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

	jsonMap := make(map[string][]entities.Rating)
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

	msg.Text = message
	return msg
}

func getHeroByDay(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	s := config.Strava{BaseUrl: os.Getenv("STRAVA_BASE_URL")}
	clubId, err := strconv.Atoi(os.Getenv("STRAVA_METRO_GROUP_ID"))
	currentTime := time.Now().Format(time.ANSIC)
	if err != nil {
		log.Panic(err)
	}
	feed := s.Feed(clubId)
	var message = "Герой дня от " + currentTime + "\n\n"
	sort.SliceStable(feed, func(i, j int) bool {
		return feed[i].Points > feed[j].Points
	})

	for i, items := range feed {
		athleteLink := fmt.Sprintf("https://www.strava.com/athletes/%v", items.AthleteId)
		message += fmt.Sprintf("%v. [%s](%s) - ",
			i+1,
			items.AthleteName,
			athleteLink)
		if items.SwimTotal > 0 {
			message += fmt.Sprintf("🏊‍♂ %.2f м, ", items.SwimTotal)
		}
		if items.BikeTotal > 0 {
			message += fmt.Sprintf("🚴 %.2f км, ", items.BikeTotal)
		}
		if items.RunTotal > 0 {
			message += fmt.Sprintf("🏃 %.2f км ⛰ %d м, ", items.RunTotal, items.ElevationGain)
		}
		message += fmt.Sprintf("*%.f ūpai* \n", items.Points)
	}
	message += "\n\n*Как начисляется ūpai за день?*\n\n"
	message += "Плавание - за 200 м плавания - 1 ūpai\n"
	message += "Вело - за 2 км езды - 1 ūpai\n"
	message += "Бег - за 1 км бега - 1 ūpai, за 100 метров подъема - 10 ūpai\n\n"
	message += "**Хотите участвовать в рейтинге дня ☀?** \n Подписывайтесь на страницу в [STRAVA](https://www.strava.com/clubs/540448)"
	msg.Text = message
	return msg
}

func replyMessage(msg tgbotapi.MessageConfig, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}
