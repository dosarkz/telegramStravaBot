package workouts

import (
	"telegramStravaBot/domain/users"
	"time"
)

type Workout struct {
	Id           int
	Title        string
	Description  string
	Status       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	WorkoutUsers []WorkoutUser
}

type WorkoutStatus struct {
	WorkoutId    int
	UserId       int64
	DeleteStatus int
	CreateStatus int
}

type WorkoutUser struct {
	Id        int
	UserID    int
	WorkoutId uint
	User      users.User
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type WorkoutUserResponse struct {
	Id        int       `json:"id"`
	UserID    int       `json:"user_id"`
	WorkoutId uint      `json:"workout_id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deleted_at"`
}

type WorkoutResponse struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ListWorkoutResponse struct {
	Data []WorkoutResponse `json:"repositories"`
}

type ListWorkoutUsersResponse struct {
	Data []WorkoutUserResponse `json:"repositories"`
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
