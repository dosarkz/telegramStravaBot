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
