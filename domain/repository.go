package domain

import (
	"github.com/jinzhu/gorm"
	"telegramStravaBot/domain/users"
	"telegramStravaBot/domain/workouts"
)

type Repositories struct {
	*users.UserRepository
	*workouts.WorkoutRepository
}

func New(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepository:    users.New(db),
		WorkoutRepository: workouts.New(db),
	}
}
