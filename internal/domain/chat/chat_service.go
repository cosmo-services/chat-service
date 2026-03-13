package chat

import (
	"errors"
	"time"
)

type ChatService struct {
	userSync    *UserProfileSync
	chatRepo    ChatRepository
	msgRepo     MessageRepository
	chatFactory *ChatFactory
	msgFactory  *MessageFactory
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

func (s *ChatService) GetDirectChat(firstUserId string, secodndUserId string) (*Chat, error) {
	direct, err := s.chatRepo.GetDirectChat(firstUserId, secodndUserId)
	if err == nil {
		return direct, nil
	}

	if !errors.Is(err, ErrChatNotFound) {
		return nil, err
	}

	direct, err = s.chatFactory.CreateDirectChat(firstUserId, secodndUserId)
	if err != nil {
		return nil, err
	}

	return direct, err
}

func (s *ChatService) SendDirectMessage(userId string, recipientUsername string, content string) error {
	if err := s.userSync.EnsureUserExists(UserSearchOptions{UserID: userId}); err != nil {
		return err
	}

	recipient, err := s.userSync.GetUser(UserSearchOptions{Username: recipientUsername})
	if err != nil {
		return err
	}

	direct, err := s.GetDirectChat(userId, recipient.UserID)
	if err != nil {
		return err
	}

	msg, err := s.msgFactory.CreateTextMessage(userId, content)
	if err != nil {
		return err
	}

	if err := s.PersistMessage(userId, direct, msg); err != nil {
		return err
	}

	return nil
}

func (s *ChatService) PersistMessage(userId string, chat *Chat, msg *Message) error {
	if !chat.IsPersisted() {
		if err := s.chatRepo.Create(chat); err != nil {
			return err
		}
	}

	if err := msg.BindToChat(chat); err != nil {
		return err
	}

	if err := s.msgRepo.CreateMessage(msg); err != nil {
		return err
	}

	if err := s.chatRepo.MarkUpdated(chat.ID, time.Now()); err != nil {
		return err
	}

	return nil
}
