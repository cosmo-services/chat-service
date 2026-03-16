package chat

import (
	"time"
)

type ChatType string

const (
	ChatTypeDirect ChatType = "direct"
	ChatTypeGroup  ChatType = "group"
)

type Chat struct {
	ID          string        `json:"id"`
	Type        ChatType      `json:"type"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	CreatedBy   string        `json:"created_by"`
	Deleted     bool          `json:"deleted"`
	Members     []*ChatMember `json:"members"`
	AvatarUrl   string        `json:"avatar_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (c *Chat) IsPersisted() bool {
	return c.ID != ""
}

func (c *Chat) GetMembersId() []string {
	memberIDs := make([]string, 0, len(c.Members))
	for _, member := range c.Members {
		if member != nil && member.UserID != "" {
			memberIDs = append(memberIDs, member.UserID)
		}
	}
	return memberIDs
}

func (c *Chat) GetAdminIds() []string {
	adminIDs := make([]string, 0)
	for _, member := range c.Members {
		if member != nil && member.UserID != "" && member.Role == RoleAdmin {
			adminIDs = append(adminIDs, member.UserID)
		}
	}
	return adminIDs
}

func (c *Chat) GetMemberIdsExcluding(excludeUserID string) []string {
	memberIDs := make([]string, 0, len(c.Members))
	for _, member := range c.Members {
		if member != nil && member.UserID != "" && member.UserID != excludeUserID {
			memberIDs = append(memberIDs, member.UserID)
		}
	}
	return memberIDs
}

func (c *Chat) GetMemberIdsByRole(role ChatMemberRole) []string {
	memberIDs := make([]string, 0)
	for _, member := range c.Members {
		if member != nil && member.UserID != "" && member.Role == role {
			memberIDs = append(memberIDs, member.UserID)
		}
	}
	return memberIDs
}

func (c *Chat) HasMember(userID string) bool {
	for _, member := range c.Members {
		if member != nil && member.UserID == userID {
			return true
		}
	}
	return false
}

func (c *Chat) GetMemberRole(userID string) (ChatMemberRole, bool) {
	for _, member := range c.Members {
		if member != nil && member.UserID == userID {
			return member.Role, true
		}
	}
	return "", false
}
