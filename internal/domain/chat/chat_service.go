package chat

import (
	"errors"
	"time"
)

const (
	NewMessageEvent string = "new_message"
)

type ChatService struct {
	userSync    *UserProfileSync
	chatRepo    ChatRepository
	chatFactory *ChatFactory
	msgRepo     MessageRepository
	msgFactory  *MessageFactory
	viewService *ViewService
	chatQuery   ChatQuery
	publisher   Publisher
}

func NewChatService(
	userSync *UserProfileSync,
	chatRepo ChatRepository,
	chatFactory *ChatFactory,
	msgRepo MessageRepository,
	msgFactory *MessageFactory,
	viewService *ViewService,
	chatQuery ChatQuery,
	publisher Publisher,
) *ChatService {
	return &ChatService{
		userSync:    userSync,
		chatRepo:    chatRepo,
		chatFactory: chatFactory,
		msgRepo:     msgRepo,
		msgFactory:  msgFactory,
		viewService: viewService,
		chatQuery:   chatQuery,
		publisher:   publisher,
	}
}

func (s *ChatService) GetDirectChatView(firstUserId string, secodnUsername string, filter MessageSearchFilter) (*CollectionView, error) {
	var directCollection *Collection
	var err error

	user1, err := s.userSync.GetUser(UserSearchOptions{UserID: firstUserId})
	if err != nil {
		return nil, err
	}

	user2, err := s.userSync.GetUser(UserSearchOptions{Username: secodnUsername})
	if err != nil {
		return nil, err
	}

	direct, err := s.GetDirectChat(user1.ID, user2.ID)
	if err != nil {
		return nil, err
	}

	if direct.IsPersisted() {
		directCollection, err = s.chatQuery.GetMessageHistory(direct.ID, filter)
		if err != nil {
			return nil, err
		}
	} else {
		directCollection = NewEmptyDirectCollection(direct, user1, user2)
	}

	directCollectionView, err := s.viewService.ProjectCollectionForViewer(firstUserId, directCollection)
	if err != nil {
		return nil, err
	}

	return directCollectionView, nil
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
	sender, err := s.userSync.GetUser(UserSearchOptions{UserID: userId})
	if err != nil {
		return err
	}

	recipient, err := s.userSync.GetUser(UserSearchOptions{Username: recipientUsername})
	if err != nil {
		return err
	}

	direct, err := s.GetDirectChat(userId, recipient.ID)
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

	collection := NewDirectMessageCollection(direct, msg, sender, recipient)
	if err := s.deliverToChatMembers(direct, collection, NewMessageEvent); err != nil {
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

func (s *ChatService) deliverToChatMembers(chat *Chat, collection *Collection, eventType string) error {
	collViewMap, err := s.viewService.ProjectCollectionForViewers(chat.GetMembersId(), collection)
	if err != nil {
		return err
	}

	if err := s.publisher.PublishToUsers(eventType, collViewMap); err != nil {
		return err
	}

	return nil
}
