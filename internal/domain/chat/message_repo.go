package chat

type MessageRepository interface {
	CreateMessage(message *Message) error
	GetMessageById(id string) (*Message, error)
	UpdateMessage(message *Message) error
	DeleteMessage(id string) error
	GetLastChatMessage(chatId string) (*Message, error)
}
