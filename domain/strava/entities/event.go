package entities

type Event struct {
	Id           int32  `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ActivityType string `json:"activity_type"`
	RouteId      int32  `json:"route_id"`
	Route        struct {
		Id   int32  `json:"id"`
		Name string `json:"name"`
		Map  struct {
			Id              string `json:"id"`
			SummaryPolyline string `json:"summary_polyline"`
		}
	}
	UpcomingOccurrences []string  `json:"upcoming_occurrences"`
	Address             string    `json:"address"`
	StartLatLng         []float64 `json:"start_latlng"`
}
