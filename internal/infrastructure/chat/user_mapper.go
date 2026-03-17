package chat_infrastructure

import (
	chat_domain "main/internal/domain/chat"

	"gorm.io/gorm"
)

func ToDomainUser(s *UserSchema) *chat_domain.User {
	if s == nil {
		return nil
	}
	return &chat_domain.User{
		ID:          s.ID,
		Username:    s.Username,
		DisplayName: s.DisplayName,
		AvatarUrl:   s.AvatarUrl,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		DeletedAt:   s.DeletedAt.Time,
	}
}

func ToSchemaUser(u *chat_domain.User) *UserSchema {
	if u == nil {
		return nil
	}
	return &UserSchema{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		AvatarUrl:   u.AvatarUrl,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   gorm.DeletedAt{Time: u.DeletedAt},
	}
}
