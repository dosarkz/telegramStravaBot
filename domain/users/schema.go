package users

import (
	"time"
)

type UserData struct {
	Id         int `gorm:"primaryKey"`
	Username   string
	TelegramId int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
