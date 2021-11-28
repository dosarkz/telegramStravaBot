package workouts

import (
	"time"
)

type Workout struct {
	Id          int `gorm:"primaryKey"`
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
