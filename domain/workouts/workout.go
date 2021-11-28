package workouts

import (
	"time"
)

type Workout struct {
	Id          int
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type WorkoutResponse struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func ToResponseModel(entity *Workout) *WorkoutResponse {
	return &WorkoutResponse{
		Id:          entity.Id,
		Title:       entity.Title,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
