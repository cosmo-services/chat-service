package chat

import "main/internal/domain/social"

func MapProfileToUser(profile *social.UserProfile) *User {
	return &User{
		UserID:      profile.UserID,
		Username:    profile.Username,
		DisplayName: profile.DisplayName,
		AvatarUrl:   profile.AvatarUrl,
	}
}
