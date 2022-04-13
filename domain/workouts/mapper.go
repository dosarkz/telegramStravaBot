package workouts

func toDBModel(entity *Workout) *Workout {
	return &Workout{
		Id:          entity.Id,
		Title:       entity.Title,
		Description: entity.Description,
	}
}

func toDomainModel(entity *Workout) *Workout {
	return &Workout{
		Id:          entity.Id,
		Title:       entity.Title,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func toWorkoutUserDBModel(entity *WorkoutUser) *WorkoutUser {
	return &WorkoutUser{
		Id:        entity.Id,
		UserID:    entity.UserID,
		WorkoutId: entity.WorkoutId,
	}
}

func toWorkoutUserDomainModel(entity *WorkoutUser) *WorkoutUser {
	return &WorkoutUser{
		Id:        entity.Id,
		UserID:    entity.UserID,
		WorkoutId: entity.WorkoutId,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
