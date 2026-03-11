package chat

type ChatQuery interface {
	GetChatHistory(chatId string, direction string, limit int) (*Collection, error)
	GetMessageHistory(chatId string, cursor string, direction string, limit int) (*Collection, error)
	GetChatWithRelations(chatId string) (*Collection, error)
	GetMessageWithRelations(messageId string) (*Collection, error)
}
