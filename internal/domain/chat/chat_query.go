package chat

type ChatQuery interface {
	GetChatHistory(userId string, filter ChatPageFilter) (*Collection, error)
	GetMessageHistory(chatId string, filter MessagePageFilter) (*Collection, error)
	GetChatWithRelations(chatId string) (*Collection, error)
	GetMessageWithRelations(messageId string) (*Collection, error)
}
