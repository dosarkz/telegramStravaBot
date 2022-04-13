package entities

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
