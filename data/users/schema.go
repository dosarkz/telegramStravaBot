package users

import (
	"time"
)

type User struct {
	Id         int `gorm:"primaryKey"`
	Username   string
	TelegramId int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
