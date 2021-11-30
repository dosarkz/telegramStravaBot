package workouts

import (
	"telegramStravaBot/data/users"
	"time"
)

type Workout struct {
	Id           int `gorm:"primaryKey"`
	Title        string
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	WorkoutUsers []WorkoutUser `gorm:"foreignKey:WorkoutId,references:Id"`
}

type WorkoutUser struct {
	Id        int `gorm:"primaryKey"`
	UserID    int
	User      users.User
	WorkoutId uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"default:null"`
}
