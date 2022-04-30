package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var strava *Strava

type Strava struct {
	BaseUrl string
}

func (s *Strava) Feed(clubId int) {
	url := fmt.Sprintf("%s/clubs/%s/feed?feed_type=club",
		s.BaseUrl,
		strconv.Itoa(clubId),
	)
	f := FeedActivity{}
	var w []WeekActivity
	Weekday := time.Now().Weekday()
	GetRequest(url, &f)

	for _, items := range f.EntriesData {
		if items.Activity.TimeAndLocation.DisplayDate == "Yesterday" && items.Activity.Type == "Run" {
			activityId, _ := strconv.ParseInt(items.Activity.Id, 10, 64)
			curWeek := fmt.Sprintf("%s/athletes/%s/goals/current_week",
				s.BaseUrl,
				items.Activity.Athlete.AthleteId,
			)
			GetRequest(curWeek, &w)
			for _, wItems := range w {
				if wItems.Sport == "run" {
					var activities []WeekItem
					switch Weekday.String() {
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
						if aItems.Id == activityId {
							fmt.Println("item:", aItems)
						}
					}
				}
			}
		}
	}
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
	CursorData struct {
		UpdatedAt int32 `json:"updated_at"`
	} `json:"cursorData"`
}

func GetRequest(url string, m any) any {
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panic(err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	req.Header.Add("x-requested-with", `XMLHttpRequest`)
	req.Header.Add("Accept", `text/javascript, application/javascript, application/ecmascript, application/x-ecmascript`)
	req.Header.Add("Cookie", `sp=ce7f9fb8-8a63-41c5-895e-5a75e94afc35; _ga=GA1.2.1213001791.1633620772; xp_session_identifier=5c74efa0e1cb60db683f77e4e1cc03cb; _strava_cbv2=true; _strava4_session=hq9ai5bthvsfqmcru17g7n11sees3hb; fbm_284597785309=base_domain=.www.strava.com; _gid=GA1.2.229978595.1651301728; CloudFront-Key-Pair-Id=APKAIDPUN4QMG7VUQPSA; _sp_id.f55d=bd9e76bf-912b-4d10-a870-569e64edb826.1651302218.1.1651302218.1651302218.325e6c2c-e943-4481-8d1d-1ab87450ae99; CloudFront-Policy=eyJTdGF0ZW1lbnQiOiBbeyJSZXNvdXJjZSI6Imh0dHBzOi8vaGVhdG1hcC1leHRlcm5hbC0qLnN0cmF2YS5jb20vKiIsIkNvbmRpdGlvbiI6eyJEYXRlTGVzc1RoYW4iOnsiQVdTOkVwb2NoVGltZSI6MTY1MjE2NjIyMH0sIkRhdGVHcmVhdGVyVGhhbiI6eyJBV1M6RXBvY2hUaW1lIjoxNjUwOTQyMjIwfX19XX0_; CloudFront-Signature=MqeRm1sUGjsD28zlEHgqd9TE3nFX8C0fhJic8YQsmUcXJZUrV95-ZN-~oAQpPMq2XbznPzGUyEPrROHahFcKdKt77LyvTD3NNAuh~VPAa56v0Sv0NLZHz6L11ITFL7-rAqzM82-YkW~eMc1ANdBUEHzvkWn3LUbylteXPYSBHbs4XGr1Wf6ZaK3zoY1SDlgb73BiSdaWl6cLW7VON4M8NwGCNT-8~Pun2-9S4YUaWKignTSloGJWqXmNf2tzcSm4kOysiNfilc~jLRn9qXT2YkNu5f5dL8YVrDjVJOE3glZmiaLDl-5LWyAITFRbPWLBoMWeK1lFNsnZufsp~t2GbQ__; _sp_ses.047d=*; fbsr_284597785309=GP2q1QaCvhHmW4v-NaL46i8HbfW7xgxtPG8qxWarM8w.eyJ1c2VyX2lkIjoiMjE0MDA3OTkwOTM4NzgzOSIsImNvZGUiOiJBUUJrcTM5NG15WXpxVnZXRHNiaFhXUmdZeDNhNTAtWXU4Uk44aTdhZHdiQkNDSmtkeVhhOXJiYk5SclA3YTN4N1lEXzJlSXQ4dTJBWTk2c1ZaR0tLbEU1VlFXT1h5RGFZMFNrSFN0N3F1Tmp5N2Vlb3d1VTlqaVRYak9HS0xBalBBSFRyT2ZSVW1iSE55UGIwZHVjUUZwc0tpOVVEaV9SWjJrSnZqQlZRUDhoS1h5TVhiNkZzbFZKVm5SMi1CVHBHcjYtN2JUb1FaSWFMS1pSb1dXMEprQWZGRF9QUGFrbS0yLURDTkF2dnIzNmVSWlhRN2ZyNDQtWlFwVkFlZXRlb2lmLUpYcG9WUlZ0WlpXcndMZUdJdnh0Q3B2OWV5UnNIT0tXTnpQcUg1NkNHMlpQTHRzUE1nMHY1YnplaE9qWk91bWZVYVBteU03X3BIQTRvMmV5VDY5UiIsIm9hdXRoX3Rva2VuIjoiRUFBQUFRa05aQWt0MEJBSEx0bGpmOWF3NHIwbE5SUzhxeG96dnNldnhaQVBmMDdWQ0ppemV0VzRUSVpCNVhSQjBhdE91UWhTQWo5NmNaQURxaTdHWkJTeVA0ZjVLWkN2clZFTG12NTIweElFYmVaQjhOWkF1NTduRHFsWVVqMjllbWlETFNkYlRnMUFBWkJHa1F0Zm5lSDJ1bFJVZ1ZVaDg2QjVvdDRnYXI5TWlhdTdFbDFrd3R0aW4xWGFlWkFjdVpDUXJ2Q1g4amJZMFpCcXByMFhWR3hiYU9vOW4iLCJhbGdvcml0aG0iOiJITUFDLVNIQTI1NiIsImlzc3VlZF9hdCI6MTY1MTMxODkxM30; _dc_gtm_UA-6309847-24=1; fbsr_284597785309=bBwxXUoKMFn0LUvZ7p-9blN2K3eceULxxTXVhboMjKo.eyJ1c2VyX2lkIjoiMjE0MDA3OTkwOTM4NzgzOSIsImNvZGUiOiJBUUR5UWVpTkNkUWwwXzctTzFQb0dqN2o5MjZEdTV6MTFteElQSnF0OFg5Z1VFMXpZNW12OU1qeWFsQ1RnQ0stMHJsUTMxSVhQMERId3VyLWxkTWV1dmJEdEFObzNpeHFMVU9XMGxpck9WaXJRLVRlUVZLM29yX0VsVlJpeWlqZ2FqT3dTVTZIV0E0TFZoc2FrdXhCR3NhaS00SkhDaWhKQzROOXRqSWhoX2hPQ2x3WjBNTHc3QzdrY2lSZUx0TmtFbC0tdnpMR0dHY0ZOX1RiVmpFS3RqY29UVVJDM0tRZUI2MTBSYmM0eWY1MkpJNFVaMjVCVWZ0Nk52enZvRGxDX1lRMGFJRFdjcGxnUXFuekEyMkttSXU0RzNWYUhuTWYwS0h6Yk5nOXNjX0dSc3JyWnZuekFiY05TRDdEWFB4c01xaW1oLW1UOWpJM1VnU2hjVUN6Y0tpcCIsIm9hdXRoX3Rva2VuIjoiRUFBQUFRa05aQWt0MEJBTlByWGJaQU5aQWlxYlNFRW55QmoweEFuWkJSSVpBS3RPa0ZkTmcwT3VjYTc5aVpBbmNkMEMxa05ha2V6d29aQ2ZGWVVZOVpCODBaQWM4NnpqMU93cThaQWxLN3BEa29MR25kN0xMQTh2N1lQaVhqQlpDU0tnWHdzbE9rWE1aQ2RFUXZoZFNINEJXUkVLY285YjgxOEZlVHFLT0VYcnJkazQ5dmEwemI4U2E5VGg0dUMyNW1SUWZPRlZuZHhuS2l5cFpDQU9QTUppeGlHWXREIiwiYWxnb3JpdGhtIjoiSE1BQy1TSEEyNTYiLCJpc3N1ZWRfYXQiOjE2NTEzMTc5NzJ9; _sp_id.047d=bf52cac3-6d67-4d51-9998-d173bcfec686.1651190399.24.1651319156.1651313561.68b49b3e-576b-4a15-8011-95d65f5609b8`)

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

	jsonErr := json.Unmarshal(body, &m)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return m
}
