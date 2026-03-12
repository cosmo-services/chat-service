package chat

type ChatService struct {
}

func NewChatService() *ChatService {
	return &ChatService{}
}

func (s *ChatService) EnsureUserExists(userId string)