package chat_infrastructure

import (
	chat_domain "main/internal/domain/chat"

	"gorm.io/gorm"
)

func ToDomainMessage(m *MessageSchema) *chat_domain.Message {
	if m == nil {
		return nil
	}

	return &chat_domain.Message{
		ID:        m.ID,
		ChatID:    m.ChatID,
		SenderID:  m.SenderID,
		Content:   m.Content,
		ReplyToId: m.ReplyToID,
		Deleted:   !m.DeletedAt.Time.IsZero(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt.Time,
	}
}

func ToSchemaMessage(m *chat_domain.Message) *MessageSchema {
	if m == nil {
		return nil
	}

	return &MessageSchema{
		ID:        m.ID,
		ChatID:    m.ChatID,
		SenderID:  m.SenderID,
		Content:   m.Content,
		ReplyToID: m.ReplyToId,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: gorm.DeletedAt{Time: m.DeletedAt},
	}
}
