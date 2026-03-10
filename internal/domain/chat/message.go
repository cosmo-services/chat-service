package chat

import "time"

type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	ReplyToId string    `json:"reply_to_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Message) BindToChat(chat *Chat) error {
	if m.ChatID != "" {
		return ErrMessageAlreadyBound
	}

	if !chat.IsPersisted() {
		return ErrChatNotFound
	}

	m.ChatID = chat.ID

	return nil
}

func (m *Message) HasReply() bool {
	return m.ReplyToId != ""
}
