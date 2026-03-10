package chat

import "time"

type ChatRepository interface {
	Create(chat *Chat) error
	GetChatById(id string) (*Chat, error)
	GetDirectChat(firstUserID, secondUserID string) (*Chat, error)
	Update(chat *Chat) error
	Delete(id string) error
	MarkUpdated(chatID string, updateTime time.Time) error
	ChatExists(chatId string) (bool, error)
	DirectChatExists(firstUserID, secondUserID string) (bool, error)
	UserInChat(userId, chatId string) (bool, error)
}
