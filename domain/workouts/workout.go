package workouts

import (
	"time"
)

type Workout struct {
	Id           int
	Title        string
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	WorkoutUsers []WorkoutUser
}

type WorkoutUser struct {
	Id        int
	UserID    uint
	WorkoutId uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type WorkoutUserResponse struct {
	Id        int       `json:"id"`
	UserID    uint      `json:"user_id"`
	WorkoutId uint      `json:"workout_id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WorkoutResponse struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ListWorkoutResponse struct {
	Data []WorkoutResponse `json:"data"`
}

type ListWorkoutUsersResponse struct {
	Data []WorkoutUserResponse `json:"data"`
}

func ToResponseWorkoutUsersModel(entity *WorkoutUser) *WorkoutUserResponse {
	return &WorkoutUserResponse{
		Id:        entity.Id,
		UserID:    entity.UserID,
		WorkoutId: entity.WorkoutId,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
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
