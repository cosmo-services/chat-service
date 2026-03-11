package chat_infrastructure

import (
	"time"

	"gorm.io/gorm"
)

type MessageSchema struct {
	ID        string         `gorm:"column:id;primaryKey"`
	ChatID    string         `gorm:"column:chat_id;not null"`
	SenderID  string         `gorm:"column:sender_id;not null"`
	Content   string         `gorm:"column:content;not null"`
	ReplyToID string         `gorm:"column:reply_to_id;index:idx_reply_to"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	Reply *MessageSchema `gorm:"foreignKey:ReplyToID;references:ID"`
	User  *UserSchema    `gorm:"foreignKey:SenderID;references:ID"`
	Chat  *ChatSchema    `gorm:"foreignKey:ChatID;references:ID"`
}

func (MessageSchema) TableName() string {
	return "messages"
}
