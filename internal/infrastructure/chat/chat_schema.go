package chat_infrastructure

import (
	"time"

	"gorm.io/gorm"
)

type ChatSchema struct {
	ID          string          `gorm:"primaryKey"`
	Type        string          `gorm:"not null"`
	Name        string          `gorm:"not null"`
	Description string          `gorm:"type:text"`
	CreatedBy   string          `gorm:"not null;index:idx_chats_created_by"`
	AvatarUrl   string          `gorm:"type:text"`
	Members     []*MemberSchema `gorm:"foreignKey:ChatID;references:ID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt  `gorm:"index"`
}

func (ChatSchema) TableName() string {
	return "chats"
}

type MemberSchema struct {
	ChatID   string    `gorm:"primaryKey;index:idx_members_chat"`
	UserID   string    `gorm:"primaryKey;index:idx_members_user"`
	Role     string    `gorm:"not null"`
	JoinedAt time.Time `gorm:"autoCreateTime"`

	User *UserSchema `gorm:"foreignKey:UserID;references:UserID"`
}

func (MemberSchema) TableName() string {
	return "chat_members"
}
