package chat

type ChatQuery interface {
	GetChatHistory(filter ChatSearchFilter) (*Collection, error)
	GetMessageHistory(chatId string, filter MessageSearchFilter) (*Collection, error)
	GetChatWithRelations(chatId string) (*Collection, error)
	GetMessageWithRelations(messageId string) (*Collection, error)
}
