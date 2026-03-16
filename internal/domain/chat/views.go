package chat

import "time"

type MessageView struct {
	ID          string    `json:"id"`
	ChatID      string    `json:"chat_id"`
	SenderID    string    `json:"sender_id"`
	Content     string    `json:"content"`
	ReplyToId   string    `json:"reply_to_id,omitempty"`
	Edited      bool      `json:"edited"`
	ContextType string    `json:"context_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type ChatView struct {
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	Name        string        `json:"name"`
	AvatarUrl   string        `json:"avatar_url"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
	Members     []*MemberView `json:"members"`
}

type MemberView struct {
	UserID   string    `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type UserView struct {
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	AvatarUrl   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type CollectionView struct {
	Chats    map[string]*ChatView    `json:"chats"`
	Messages map[string]*MessageView `json:"messages"`
	Users    map[string]*UserView    `json:"users"`

	HasNext bool `json:"has_next"`
	HasPrev bool `json:"has_prev"`
}
