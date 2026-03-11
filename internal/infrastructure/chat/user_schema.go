package chat_infrastructure

import (
	"time"

	"gorm.io/gorm"
)

type UserSchema struct {
	ID          string         `gorm:"primaryKey"`
	Username    string         `gorm:"not null;uniqueIndex:idx_users_username"`
	DisplayName string         `gorm:"not null"`
	AvatarUrl   string         `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (UserSchema) TableName() string {
	return "users"
}
