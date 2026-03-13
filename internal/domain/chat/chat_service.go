package chat

type ChatService struct {
	userSync    *UserProfileSync
	chatRepo    ChatRepository
	chatFactory *ChatFactory
	userRepo    UserRepository
}

func NewChatService(
	userSync *UserProfileSync,
	chatRepo ChatRepository,
) *ChatService {
	return &ChatService{
		userSync: userSync,
		chatRepo: chatRepo,
	}
}

func (s *ChatService) GetDirectChat(userId string, companionUsername string) (*Chat, error) {
	companion, err := s.userSync.GetUser(UserSearchOptions{Username: companionUsername})
	if err != nil {
		return nil, err
	}

	direct, err := s.chatRepo.GetDirectChat(userId, companion.UserID)
	if err != nil {
		return nil, err
	}

	if direct == nil {
		direct, err = s.chatFactory.CreateDirectChat(userId, companion.UserID)
		if err != nil {
			return nil, err
		}
	}

	return direct, err
}
