package users

import (
	domain "telegramStravaBot/domain/users"
)

func toDBModel(entity *domain.User) *User {
	return &User{
		Id:         entity.Id,
		Username:   entity.Username,
		TelegramId: entity.TelegramId,
	}
}

func toDomainModel(entity *User) *domain.User {
	return &domain.User{
		Id:         entity.Id,
		Username:   entity.Username,
		TelegramId: entity.TelegramId,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
}
