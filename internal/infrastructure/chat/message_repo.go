package chat_infrastructure

import (
	"errors"
	chat_domain "main/internal/domain/chat"

	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) chat_domain.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) CreateMessage(message *chat_domain.Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}

	schema := ToSchemaMessage(message)

	return r.db.Create(schema).Error
}

func (r *messageRepository) GetMessageById(id string) (*chat_domain.Message, error) {
	var msgSchema MessageSchema

	err := r.db.First(&msgSchema, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, chat_domain.ErrMessageNotFound
		}
		return nil, err
	}

	return ToDomainMessage(&msgSchema), nil
}

func (r *messageRepository) UpdateMessage(message *chat_domain.Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}

	schema := ToSchemaMessage(message)

	result := r.db.Save(schema)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return chat_domain.ErrMessageNotFound
	}

	return nil
}

func (r *messageRepository) DeleteMessage(id string) error {
	result := r.db.Delete(&MessageSchema{}, "id = ?", id)

	if result.RowsAffected == 0 {
		return chat_domain.ErrMessageNotFound
	}

	return result.Error
}

func (r *messageRepository) GetLastChatMessage(chatID string) (*chat_domain.Message, error) {
	var msgSchema MessageSchema

	err := r.db.
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		First(&msgSchema).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, chat_domain.ErrMessageNotFound
		}
		return nil, err
	}

	return ToDomainMessage(&msgSchema), nil
}

func (r *messageRepository) GetMessageByIdUnscoped(id string) (*chat_domain.Message, error) {
	var msgSchema MessageSchema

	err := r.db.
		Unscoped().
		First(&msgSchema, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, chat_domain.ErrMessageNotFound
		}
		return nil, err
	}

	return ToDomainMessage(&msgSchema), nil
}

func (r *messageRepository) RestoreMessage(id string) error {
	err := r.db.
		Unscoped().
		Model(&MessageSchema{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error

	if err != nil {
		return err
	}
	return nil
}
