package chat

import (
	"errors"
	"time"
)

const (
	NewMessageEvent string = "new_message"
	NewChatEvent    string = "new_chat"
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

func (s *ChatService) GetChatsHistoryView(userId string, filter ChatPageFilter) (*CollectionView, error) {
	if err := s.userSync.EnsureUserExists(UserSearchOptions{UserID: userId}); err != nil {
		return nil, err
	}

	collection, err := s.chatQuery.GetChatHistory(filter)
	if err != nil {
		return nil, err
	}

	collcetionView, err := s.viewService.ProjectCollectionForViewer(userId, collection)
	if err != nil {
		return nil, err
	}

	return collcetionView, nil
}

func (s *ChatService) GetDirectMessageHistoryView(firstUserId string, secondUsername string, filter MessagePageFilter) (*CollectionView, error) {
	var directCollection *Collection
	var err error

	user1, err := s.userSync.GetUser(UserSearchOptions{UserID: firstUserId})
	if err != nil {
		return nil, err
	}

	user2, err := s.userSync.GetUser(UserSearchOptions{Username: secondUsername})
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

func (s *ChatService) GetChatMessageHistoryView(userId string, chatId string, filter MessagePageFilter) (*CollectionView, error) {
	collection, err := s.chatQuery.GetMessageHistory(chatId, filter)
	if err != nil {
		return nil, err
	}

	collectionView, err := s.viewService.ProjectCollectionForViewer(userId, collection)
	if err != nil {
		return nil, err
	}

	return collectionView, nil
}

func (s *ChatService) GetDirectChat(firstUserId string, secondUserId string) (*Chat, error) {
	direct, err := s.chatRepo.GetDirectChat(firstUserId, secondUserId)
	if err == nil {
		return direct, nil
	}

	if !errors.Is(err, ErrChatNotFound) {
		return nil, err
	}

	direct, err = s.chatFactory.CreateDirectChat(firstUserId, secondUserId)
	if err != nil {
		return nil, err
	}

	return direct, err
}

func (s *ChatService) GetChatForUser(userId string, chatId string) (*Chat, error) {
	chat, err := s.chatRepo.GetChatById(chatId)
	if err != nil {
		return nil, err
	}

	if !chat.HasMember(userId) {
		return nil, ErrChatAccessDenied
	}

	return chat, nil
}

func (s *ChatService) SendChatMessage(userId string, chatId string, content string) error {
	var collection *Collection
	var err error

	chat, err := s.GetChatForUser(userId, chatId)
	if err != nil {
		return err
	}

	msg, err := s.msgFactory.CreateTextMessage(userId, content)
	if err != nil {
		return err
	}

	if err := s.persistMessage(chat, msg); err != nil {
		return err
	}

	sender, err := s.userSync.GetUser(UserSearchOptions{UserID: userId})
	if err != nil {
		return err
	}

	if chat.Type == ChatTypeDirect {
		companionId := chat.GetMemberIdsExcluding(userId)[0]
		companion, err := s.userSync.GetUser(UserSearchOptions{UserID: companionId})
		if err != nil {
			return err
		}

		collection = NewDirectMessageCollection(chat, msg, sender, companion)
	} else {
		collection = NewChatMessageCollection(chat, msg, sender)
	}

	if err := s.deliverToChatMembers(chat, collection, NewMessageEvent); err != nil {
		return err
	}

	return nil
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

	if err := s.persistMessage(direct, msg); err != nil {
		return err
	}

	collection := NewDirectMessageCollection(direct, msg, sender, recipient)
	if err := s.deliverToChatMembers(direct, collection, NewMessageEvent); err != nil {
		return err
	}

	return nil
}

func (s *ChatService) persistMessage(chat *Chat, msg *Message) error {
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
