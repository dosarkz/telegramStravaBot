package workouts

import (
	"time"
)

type Workout struct {
	Id           int `gorm:"primaryKey"`
	Title        string
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	WorkoutUsers []WorkoutUser
}

type WorkoutUser struct {
	Id        int `gorm:"primaryKey"`
	UserID    uint
	WorkoutId uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
