package yandex

type Weather struct {
	Now  int64 `json:"now"`
	Fact struct {
		Temp      int32   `json:"temp"`
		Feel      int32   `json:"feels_like"`
		Icon      string  `json:"icon"`
		Condition string  `json:"condition"`
		WindSpeed float32 `json:"wind_speed"`
		Humidity  int     `json:"humidity"`
	}
	Forecast struct {
		Sunrise string `json:"sunrise"`
	}
}
