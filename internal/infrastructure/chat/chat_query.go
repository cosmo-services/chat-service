package chat_infrastructure

import (
	"errors"

	chat_domain "main/internal/domain/chat"
	"main/pkg"

	"gorm.io/gorm"
)

type chatQuery struct {
	db pkg.GormDB
}

func NewChatQuery(db pkg.GormDB) chat_domain.ChatQuery {
	return &chatQuery{db: db}
}

func (q *chatQuery) GetChatHistory(filter chat_domain.ChatSearchFilter) (*chat_domain.Collection, error) {
	collection := &chat_domain.Collection{
		Chats:    make(map[string]*chat_domain.Chat),
		Messages: make(map[string]*chat_domain.Message),
		Replies:  make(map[string]*chat_domain.Message),
		Users:    make(map[string]*chat_domain.User),
	}

	var chats []ChatSchema

	query := q.db.DB.
		Preload("Members.User")

	page, err := pkg.Paginate(query, &chats, filter.Cursor, "updated_at", filter.Direction, filter.Limit)
	if err != nil {
		return nil, err
	}

	collection.HasNext = page.HasNext
	collection.HasPrev = page.HasPrev

	if len(chats) == 0 {
		return collection, nil
	}

	chatIDs := make([]string, len(chats))
	for i, chat := range chats {
		chatIDs[i] = chat.ID
	}

	var lastMessages []*MessageSchema
	err = q.db.DB.
		Where("chat_id IN ?", chatIDs).
		Where("(chat_id, created_at) IN (?)",
			q.db.DB.Table("messages").
				Select("chat_id, MAX(created_at)").
				Where("chat_id IN ?", chatIDs).
				Group("chat_id"),
		).
		Preload("Reply").
		Find(&lastMessages).Error
	if err != nil {
		return nil, err
	}

	for _, msg := range lastMessages {
		if msg == nil {
			continue
		}

		collection.Messages[msg.ID] = ToDomainMessage(msg)

		if msg.Reply != nil {
			if _, exists := collection.Messages[msg.Reply.ID]; !exists {
				collection.Replies[msg.Reply.ID] = ToDomainMessage(msg.Reply)
			}
		}
	}

	for _, chatSchema := range chats {
		chat := ToDomainChat(&chatSchema)
		chat.Members = make([]*chat_domain.ChatMember, 0, len(chatSchema.Members))

		for _, memberSchema := range chatSchema.Members {
			member := ToDomainMember(memberSchema)
			chat.Members = append(chat.Members, member)

			if memberSchema.User != nil {
				if _, exists := collection.Users[memberSchema.User.ID]; !exists {
					collection.Users[memberSchema.User.ID] = ToDomainUser(memberSchema.User)
				}
			}
		}

		collection.Chats[chat.ID] = chat
	}

	return collection, nil
}

func (q *chatQuery) GetMessageHistory(chatId string, filter chat_domain.MessageSearchFilter) (*chat_domain.Collection, error) {
	collection := &chat_domain.Collection{
		Chats:    make(map[string]*chat_domain.Chat),
		Messages: make(map[string]*chat_domain.Message),
		Replies:  make(map[string]*chat_domain.Message),
		Users:    make(map[string]*chat_domain.User),
	}

	var chat ChatSchema
	if err := q.db.DB.
		Preload("Members.User").
		First(&chat, "id = ?", chatId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, chat_domain.ErrChatNotFound
		}
		return nil, err
	}

	domainChat := ToDomainChat(&chat)
	domainChat.Members = make([]*chat_domain.ChatMember, 0, len(chat.Members))
	for _, memberSchema := range chat.Members {
		member := ToDomainMember(memberSchema)
		domainChat.Members = append(domainChat.Members, member)

		if memberSchema.User != nil {
			collection.Users[memberSchema.User.ID] = ToDomainUser(memberSchema.User)
		}
	}
	collection.Chats[chat.ID] = domainChat

	var messages []*MessageSchema

	query := q.db.DB.
		Where("chat_id = ?", chatId).
		Preload("Reply")

	page, err := pkg.Paginate(
		query,
		&messages,
		filter.Cursor,
		"created_at",
		filter.Direction,
		filter.Limit,
	)
	if err != nil {
		return nil, err
	}

	collection.HasNext = page.HasNext
	collection.HasPrev = page.HasPrev

	if len(messages) == 0 {
		return collection, nil
	}

	for _, msg := range messages {
		if msg == nil {
			continue
		}

		collection.Messages[msg.ID] = ToDomainMessage(msg)

		if msg.Reply != nil {
			if _, exists := collection.Messages[msg.Reply.ID]; !exists {
				collection.Replies[msg.Reply.ID] = ToDomainMessage(msg.Reply)
			}
		}
	}

	return collection, nil
}

func (q *chatQuery) GetChatWithRelations(chatId string) (*chat_domain.Collection, error) {
	collection := &chat_domain.Collection{
		Chats:    make(map[string]*chat_domain.Chat),
		Messages: make(map[string]*chat_domain.Message),
		Replies:  make(map[string]*chat_domain.Message),
		Users:    make(map[string]*chat_domain.User),
	}

	var chat ChatSchema
	if err := q.db.DB.
		Preload("Members.User").
		First(&chat, "id = ?", chatId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, chat_domain.ErrChatNotFound
		}
		return nil, err
	}

	domainChat := ToDomainChat(&chat)
	domainChat.Members = make([]*chat_domain.ChatMember, 0, len(chat.Members))
	for _, memberSchema := range chat.Members {
		member := ToDomainMember(memberSchema)
		domainChat.Members = append(domainChat.Members, member)

		if memberSchema.User != nil {
			collection.Users[memberSchema.User.ID] = ToDomainUser(memberSchema.User)
		}
	}
	collection.Chats[chat.ID] = domainChat

	var lastMsg MessageSchema
	err := q.db.DB.
		Where("chat_id = ?", chatId).
		Order("created_at DESC, id DESC").
		Preload("Reply").
		First(&lastMsg).Error

	if err == nil {
		collection.Messages[lastMsg.ID] = ToDomainMessage(&lastMsg)
		if lastMsg.Reply != nil {
			collection.Replies[lastMsg.Reply.ID] = ToDomainMessage(lastMsg.Reply)
		}
	}

	return collection, nil
}

func (q *chatQuery) GetMessageWithRelations(messageId string) (*chat_domain.Collection, error) {
	collection := &chat_domain.Collection{
		Chats:    make(map[string]*chat_domain.Chat),
		Messages: make(map[string]*chat_domain.Message),
		Replies:  make(map[string]*chat_domain.Message),
		Users:    make(map[string]*chat_domain.User),
	}

	var msg MessageSchema
	if err := q.db.DB.
		Preload("Reply").
		Preload("Chat.Members.User").
		First(&msg, "id = ?", messageId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, chat_domain.ErrMessageNotFound
		}
		return nil, err
	}

	collection.Messages[msg.ID] = ToDomainMessage(&msg)

	if msg.Reply != nil {
		collection.Replies[msg.Reply.ID] = ToDomainMessage(msg.Reply)
	}

	if msg.Chat != nil {
		chat := msg.Chat
		domainChat := ToDomainChat(chat)
		domainChat.Members = make([]*chat_domain.ChatMember, 0, len(chat.Members))

		for _, memberSchema := range chat.Members {
			member := ToDomainMember(memberSchema)
			domainChat.Members = append(domainChat.Members, member)

			if memberSchema.User != nil {
				collection.Users[memberSchema.User.ID] = ToDomainUser(memberSchema.User)
			}
		}
		collection.Chats[chat.ID] = domainChat
	}

	return collection, nil
}
