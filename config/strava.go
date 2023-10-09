package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var strava *Strava

type Strava struct {
	BaseUrl string
}

type Feed struct {
	AvatarUrl     string
	AthleteName   string
	AthleteId     string
	SwimTotal     float32
	BikeTotal     float32
	RunTotal      float32
	Points        float32
	ElevationGain int32
}

func FeedRequest(baseUrl string, clubId string, cursor int32) string {
	url := fmt.Sprintf("%s/clubs/%s/feed?feed_type=club",
		baseUrl, clubId)
	if cursor != 0 {
		url += "&cursor=" + strconv.Itoa(int(cursor))
	}

	return url
}

func (s *Strava) Feed(clubId int) []Feed {
	var feed []Feed
	url := FeedRequest(s.BaseUrl, strconv.Itoa(clubId), 0)

	f := FeedActivity{}
	var w []WeekActivity
	Weekday := time.Now().Weekday()

	//// List past user activities from club
	GetFeed(url, &f, s.BaseUrl, strconv.Itoa(clubId), &feed, w, Weekday)
	return feed
}

func GetFeed(url string, f *FeedActivity, baseUrl string, clubId string,
	feed *[]Feed, w []WeekActivity, weekday time.Weekday) *[]Feed {

	entry := make(chan []Entries)

	go GetRequest(url, entry, f)

	for e, items := range <-entry {

		found := false
		athlete := Feed{}
		athlete.AthleteId = items.Activity.Athlete.AthleteId
		athlete.AthleteName = items.Activity.Athlete.AthleteName

		for _, value := range *feed {
			if value.AthleteId == items.Activity.Athlete.AthleteId {
				found = true
				break
			}
		}

		if found {
			continue
		}

		currentDay := items.Activity.TimeAndLocation.DisplayDate

		if len(items.RowData.Activities) > 0 {
			for _, activityItem := range items.RowData.Activities {
				_, m, d := activityItem.StartDateLocal.Date()
				_, nm, nd := time.Now().Date()

				if m == nm && d == nd {
					for j, value := range *feed {
						if value.AthleteId == strconv.Itoa(int(activityItem.AthleteId)) {
							athlete = value
							break
						}

						if len(*feed) == j+1 && value.AthleteId != strconv.Itoa(int(activityItem.AthleteId)) {
							athlete = Feed{
								AthleteName: activityItem.AthleteName,
								AthleteId:   strconv.Itoa(int(activityItem.AthleteId)),
							}

							getCurrentWeekActivities(baseUrl, athlete, w, weekday, feed)
							break
						}
					}
				}

			}
		}

		if athlete.AthleteId == "" {
			continue
		}

		if e+1 == len(f.EntriesData) {
			url = FeedRequest(baseUrl, clubId, items.CursorData.UpdatedAt)
			feed = GetFeed(url, f, baseUrl, clubId, feed, w, weekday)
		}

		if currentDay == "Yesterday" {
			feed = getCurrentWeekActivities(baseUrl, athlete, w, weekday, feed)
		}

	}
	return feed
}

func getCurrentWeekActivities(baseUrl string, athlete Feed, w []WeekActivity, weekday time.Weekday, feed *[]Feed) *[]Feed {
	curWeek := fmt.Sprintf("%s/athletes/%s/goals/current_week",
		baseUrl,
		athlete.AthleteId,
	)
	we := make(chan []WeekActivity)
	// Activities by athlete from current week
	go GetCurrWeekRequest(curWeek, we, w)
	fmt.Println("Previous RUN TOTAL", athlete.RunTotal)
	for _, wItems := range <-we {

		var activities []WeekItem
		switch weekday.String() {
		case "Monday":
			activities = wItems.ByDayOfWeek.Monday.Activities
			break
		case "Tuesday":
			activities = wItems.ByDayOfWeek.Tuesday.Activities
			break
		case "Wednesday":
			activities = wItems.ByDayOfWeek.Wednesday.Activities
			break
		case "Thursday":
			activities = wItems.ByDayOfWeek.Thursday.Activities
			break
		case "Friday":
			activities = wItems.ByDayOfWeek.Friday.Activities
			break
		case "Saturday":
			activities = wItems.ByDayOfWeek.Saturday.Activities
			break
		case "Sunday":
			activities = wItems.ByDayOfWeek.Sunday.Activities
			break
		}

		for _, aItems := range activities {
			fmt.Println("Find item", aItems)
			switch aItems.Type {
			case "Swim":
				fmt.Println("Swim", aItems.Distance)
				if aItems.Distance >= 100 {
					athlete.SwimTotal += aItems.Distance
					athlete.Points += athlete.SwimTotal / 200
				}
				break
			case "Ride":
				if aItems.Distance >= 100 {
					athlete.BikeTotal += aItems.Distance / 1000
					athlete.Points += athlete.BikeTotal / 2
				}
				break
			case "Run":
				if aItems.Distance >= 100 {

					athlete.RunTotal += aItems.Distance / 1000
					fmt.Println("RUN POINTS BEFORE", athlete.Points)
					fmt.Println("RUN TOTAL", athlete.RunTotal)
					athlete.ElevationGain += aItems.ElevGain
					athlete.Points += aItems.Distance/1000 + float32(aItems.ElevGain/10)

					fmt.Println("RUN POINTS AFTER", athlete.Points)
				}
				break
			}
		}
	}

	if athlete.Points == 0 {
		return feed
	}

	for s, value := range *feed {
		if value.AthleteId == athlete.AthleteId {
			break
		}

		if len(*feed) == s+1 && value.AthleteId != athlete.AthleteId {
			*feed = append(*feed, athlete)
			break
		}
	}

	if len(*feed) == 0 {
		*feed = append(*feed, athlete)
	}

	return feed
}

type WeekActivity struct {
	Id          string `json:"id"`
	Sport       string `json:"sport"`
	ByDayOfWeek struct {
		Monday    WeekItemScope `json:"1"`
		Tuesday   WeekItemScope `json:"2"`
		Wednesday WeekItemScope `json:"3"`
		Thursday  WeekItemScope `json:"4"`
		Friday    WeekItemScope `json:"5"`
		Saturday  WeekItemScope `json:"6"`
		Sunday    WeekItemScope `json:"7"`
	} `json:"by_day_of_week"`
}

type WeekItemScope struct {
	Activities []WeekItem `json:"activities"`
}

type WeekItem struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Distance float32 `json:"distance"`
	ElevGain int32   `json:"elev_gain"`
	Type     string  `json:"type"`
	Speed    float32 `json:"speed"`
}

func GetStrava() *Strava {
	return strava
}

type FeedActivity struct {
	EntriesData []Entries `json:"entries"`
}

type Entries struct {
	Activity struct {
		Id           string `json:"id"`
		ActivityName string `json:"activityName"`
		Type         string `json:"type"`
		Athlete      struct {
			AvatarUrl   string `json:"avatarUrl"`
			AthleteName string `json:"athleteName"`
			AthleteId   string `json:"athleteId"`
			Sex         string `json:"sex"`
		} `json:"athlete"`
		TimeAndLocation struct {
			DisplayDateAtTime string `json:"displayDateAtTime"`
			DisplayDate       string `json:"displayDate"`
			Location          string `json:"location"`
		} `json:"timeAndLocation"`
	} `json:"activity"`
	RowData struct {
		Entity     string `json:"entity"`
		Activities []GroupActivityItem
	}
	CursorData struct {
		UpdatedAt int32 `json:"updated_at"`
	} `json:"cursorData"`
}

type GroupActivityItem struct {
	Entity         string      `json:"entity"`
	EntityId       int64       `json:"entity_id"`
	AthleteId      int64       `json:"athlete_id"`
	AthleteName    string      `json:"athlete_name"`
	Type           string      `json:"type"`
	StartDateLocal time.Time   `json:"start_date_local"`
	Visibility     string      `json:"visibility"`
	MapPolyline    [][]float64 `json:"map_polyline"`
	Location       string      `json:"location"`
}

func GetRequest(url string, entry chan []Entries, f *FeedActivity) {
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panic(err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	req.Header.Add("x-requested-with", `XMLHttpRequest`)
	req.Header.Add("Accept", `text/javascript, application/javascript, application/ecmascript, application/x-ecmascript`)
	req.Header.Add("Cookie", `_sp_ses.047d=*; sp=19b21ab6-e548-49a4-95ba-1500c64e431e; _strava_cbv3=true; _gid=GA1.2.1211933196.1696870983; _dc_gtm_UA-6309847-24=1; _ga=GA1.1.777890566.1696870983; _scid=0baa2b72-46bd-4f43-89c6-e05ae20c5bdf; _scid_r=0baa2b72-46bd-4f43-89c6-e05ae20c5bdf; _sctr=1%7C1696788000000; _sp_id.047d=3c059e1b-9b66-4a41-be03-44195dbef1e1.1696870981.1.1696870987.1696870981.95d56064-6830-43c8-946f-05e5dc36184a; _gcl_au=1.1.616221630.1696870983.127475789.1696870984.1696870987; _strava4_session=dkh525jctbmigaj3omc8c5ueu5k9nruq; _ga_ESZ0QKJW56=GS1.1.1696870983.1.0.1696870989.54.0.0`)

	resp, err := client.Do(req)

	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Panic(err)
			}
		}(resp.Body)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &f)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	entry <- f.EntriesData
}
func GetCurrWeekRequest(url string, wa chan []WeekActivity, m []WeekActivity) {
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panic(err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	req.Header.Add("x-requested-with", `XMLHttpRequest`)
	req.Header.Add("Accept", `text/javascript, application/javascript, application/ecmascript, application/x-ecmascript`)
	req.Header.Add("Cookie", `_sp_ses.047d=*; sp=19b21ab6-e548-49a4-95ba-1500c64e431e; _strava_cbv3=true; _gid=GA1.2.1211933196.1696870983; _dc_gtm_UA-6309847-24=1; _ga=GA1.1.777890566.1696870983; _scid=0baa2b72-46bd-4f43-89c6-e05ae20c5bdf; _scid_r=0baa2b72-46bd-4f43-89c6-e05ae20c5bdf; _sctr=1%7C1696788000000; _sp_id.047d=3c059e1b-9b66-4a41-be03-44195dbef1e1.1696870981.1.1696870987.1696870981.95d56064-6830-43c8-946f-05e5dc36184a; _gcl_au=1.1.616221630.1696870983.127475789.1696870984.1696870987; _strava4_session=dkh525jctbmigaj3omc8c5ueu5k9nruq; _ga_ESZ0QKJW56=GS1.1.1696870983.1.0.1696870989.54.0.0`)

	resp, err := client.Do(req)

	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Panic(err)
			}
		}(resp.Body)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &m)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	wa <- m
}
