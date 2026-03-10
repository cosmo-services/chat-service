package chat

import (
	"time"
)

type ChatMemberRole string

const (
	RoleMember ChatMemberRole = "member"
	RoleAdmin  ChatMemberRole = "admin"
)

type ChatMember struct {
	UserID   string         `json:"user_id"`
	Role     ChatMemberRole `json:"role"`
	JoinedAt time.Time      `json:"joined_at"`
}
