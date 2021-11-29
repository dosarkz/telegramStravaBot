package workouts

import (
	domain "telegramStravaBot/domain/workouts"
)

func toDBModel(entity *domain.Workout) *Workout {
	return &Workout{
		Id:          entity.Id,
		Title:       entity.Title,
		Description: entity.Description,
	}
}

func toDomainModel(entity *Workout) *domain.Workout {
	return &domain.Workout{
		Id:          entity.Id,
		Title:       entity.Title,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func toWorkoutUserDBModel(entity *domain.WorkoutUser) *WorkoutUser {
	return &WorkoutUser{
		Id:        entity.Id,
		UserID:    entity.UserID,
		WorkoutId: entity.WorkoutId,
	}
}

func toWorkoutUserDomainModel(entity *WorkoutUser) *domain.WorkoutUser {
	return &domain.WorkoutUser{
		Id:        entity.Id,
		UserID:    entity.UserID,
		WorkoutId: entity.WorkoutId,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
