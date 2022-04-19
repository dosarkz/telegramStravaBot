package users

func toDBModel(entity *User) *User {
	return &User{
		Id:         entity.Id,
		Username:   entity.Username,
		TelegramId: entity.TelegramId,
	}
}

func toDomainModel(entity *User) *User {
	return &User{
		Id:         entity.Id,
		Username:   entity.Username,
		TelegramId: entity.TelegramId,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
}
