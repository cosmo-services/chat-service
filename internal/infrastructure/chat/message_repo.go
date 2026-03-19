package chat_infrastructure

import (
	"errors"
	"fmt"
	chat_domain "main/internal/domain/chat"
	"main/pkg"
	"time"

	"gorm.io/gorm"
)

type messageRepository struct {
	db pkg.GormDB
}

func NewMessageRepository(db pkg.GormDB) chat_domain.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) CreateMessage(message *chat_domain.Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}

	message.ID = generateMessageID()

	schema := ToSchemaMessage(message)

	return r.db.DB.Create(schema).Error
}

func (r *messageRepository) GetMessageById(id string) (*chat_domain.Message, error) {
	var msgSchema MessageSchema

	err := r.db.DB.First(&msgSchema, "id = ?", id).Error
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

	result := r.db.DB.Save(schema)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return chat_domain.ErrMessageNotFound
	}

	return nil
}

func (r *messageRepository) DeleteMessage(id string) error {
	result := r.db.DB.Delete(&MessageSchema{}, "id = ?", id)

	if result.RowsAffected == 0 {
		return chat_domain.ErrMessageNotFound
	}

	return result.Error
}

func (r *messageRepository) GetLastChatMessage(chatID string) (*chat_domain.Message, error) {
	var msgSchema MessageSchema

	err := r.db.DB.
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

	err := r.db.DB.
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
	err := r.db.DB.
		Unscoped().
		Model(&MessageSchema{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error

	if err != nil {
		return err
	}
	return nil
}

func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}
