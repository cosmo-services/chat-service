package chat_infrastructure

import (
	"context"
	"errors"
	"time"

	chat_domain "main/internal/domain/chat"
	"main/pkg"

	"gorm.io/gorm"
)

type chatRepository struct {
	db pkg.GormDB
}

func NewChatRepository(db pkg.GormDB) chat_domain.ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) Create(chat *chat_domain.Chat) error {
	if chat == nil {
		return errors.New("chat cannot be nil")
	}

	return r.db.WithTransaction(context.Background(), func(tx *gorm.DB) error {
		schema := ToSchemaChat(chat)

		if err := tx.Create(schema).Error; err != nil {
			return err
		}

		if len(chat.Members) > 0 {
			for _, member := range chat.Members {
				memberSchema := ToSchemaMember(member)
				memberSchema.ChatID = chat.ID

				if err := tx.Create(memberSchema).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *chatRepository) Update(chat *chat_domain.Chat) error {
	if chat == nil {
		return errors.New("chat cannot be nil")
	}

	return r.db.WithTransaction(context.Background(), func(tx *gorm.DB) error {
		schema := ToSchemaChat(chat)

		if err := tx.Model(&ChatSchema{}).
			Where("id = ?", chat.ID).
			Updates(map[string]interface{}{
				"type":        schema.Type,
				"name":        schema.Name,
				"description": schema.Description,
				"updated_at":  schema.UpdatedAt,
			}).Error; err != nil {
			return err
		}
		if err := r.syncMembers(tx, chat.ID, chat.Members); err != nil {
			return err
		}

		return nil
	})
}

func (r *chatRepository) GetChatById(id string) (*chat_domain.Chat, error) {
	var chatSchema ChatSchema

	err := r.db.DB.
		Preload("Members").
		First(&chatSchema, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, chat_domain.ErrChatNotFound
		}
		return nil, err
	}

	return ToDomainChat(&chatSchema), nil
}

func (r *chatRepository) Delete(id string) error {
	result := r.db.DB.Delete(&ChatSchema{}, "id = ?", id)

	if result.RowsAffected == 0 {
		return chat_domain.ErrChatNotFound
	}

	return result.Error
}

func (r *chatRepository) GetDirectChat(firstUserID, secondUserID string) (*chat_domain.Chat, error) {
	var chatSchema ChatSchema

	err := r.db.DB.
		Preload("Members").
		Where("type = ?", chat_domain.ChatTypeDirect).
		Joins(`JOIN chat_members cm1 ON cm1.chat_id = chats.id 
			   AND cm1.user_id = ? AND cm1.deleted_at IS NULL`, firstUserID).
		Joins(`JOIN chat_members cm2 ON cm2.chat_id = chats.id 
			   AND cm2.user_id = ? AND cm2.deleted_at IS NULL`, secondUserID).
		First(&chatSchema).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return ToDomainChat(&chatSchema), nil
}

func (r *chatRepository) DirectChatExists(firstUserID, secondUserID string) (bool, error) {
	var count int64

	err := r.db.DB.
		Model(&ChatSchema{}).
		Where("type = ?", chat_domain.ChatTypeDirect).
		Joins(`JOIN chat_members cm1 ON cm1.chat_id = chats.id 
			   AND cm1.user_id = ? AND cm1.deleted_at IS NULL`, firstUserID).
		Joins(`JOIN chat_members cm2 ON cm2.chat_id = chats.id 
			   AND cm2.user_id = ? AND cm2.deleted_at IS NULL`, secondUserID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *chatRepository) ChatExists(chatID string) (bool, error) {
	var count int64

	err := r.db.DB.
		Model(&ChatSchema{}).
		Where("id = ?", chatID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *chatRepository) UserInChat(userID, chatID string) (bool, error) {
	var count int64

	err := r.db.DB.
		Model(&MemberSchema{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *chatRepository) MarkUpdated(chatID string, updateTime time.Time) error {
	err := r.db.DB.
		Model(&ChatSchema{}).
		Where("id = ?", chatID).
		Update("updated_at", updateTime).Error

	return err
}

func (r *chatRepository) GetChatByIdUnscoped(id string) (*chat_domain.Chat, error) {
	var chatSchema ChatSchema

	err := r.db.DB.
		Unscoped().
		Preload("Members").
		First(&chatSchema, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, chat_domain.ErrChatNotFound
		}
		return nil, err
	}

	return ToDomainChat(&chatSchema), nil
}

func (r *chatRepository) RestoreChat(id string) error {
	err := r.db.DB.
		Unscoped().
		Model(&ChatSchema{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error

	return err
}

func (r *chatRepository) AddMember(chatID string, member *chat_domain.ChatMember) error {
	if member == nil {
		return errors.New("member cannot be nil")
	}

	schema := ToSchemaMember(member)
	schema.ChatID = chatID

	return r.db.DB.Create(schema).Error
}

func (r *chatRepository) RemoveMember(chatID, userID string) error {
	result := r.db.DB.
		Delete(&MemberSchema{}, "chat_id = ? AND user_id = ?", chatID, userID)

	if result.RowsAffected == 0 {
		return chat_domain.ErrMemberNotFound
	}

	return result.Error
}

func (r *chatRepository) GetMembers(chatID string) ([]*chat_domain.ChatMember, error) {
	var members []MemberSchema

	err := r.db.DB.
		Where("chat_id = ?", chatID).
		Find(&members).Error

	if err != nil {
		return nil, err
	}

	result := make([]*chat_domain.ChatMember, 0, len(members))
	for _, m := range members {
		result = append(result, ToDomainMember(&m))
	}

	return result, nil
}

func (r *chatRepository) syncMembers(
	tx *gorm.DB, chatID string, newMembers []*chat_domain.ChatMember) error {
	var existingMembers []MemberSchema
	if err := tx.Where("chat_id = ?", chatID).Find(&existingMembers).Error; err != nil {
		return err
	}

	newMembersMap := make(map[string]*chat_domain.ChatMember, len(newMembers))
	for _, m := range newMembers {
		if m != nil {
			newMembersMap[m.UserID] = m
		}
	}

	existingMembersMap := make(map[string]*MemberSchema, len(existingMembers))
	for _, m := range existingMembers {
		existingMembersMap[m.UserID] = &m
	}

	for userID, newMember := range newMembersMap {
		memberSchema := ToSchemaMember(newMember)
		memberSchema.ChatID = chatID

		if _, exists := existingMembersMap[userID]; exists {
			if err := tx.Model(&MemberSchema{}).
				Where("chat_id = ? AND user_id = ?", chatID, userID).
				Update("role", memberSchema.Role).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Create(memberSchema).Error; err != nil {
				return err
			}
		}
	}

	for userID := range existingMembersMap {
		if _, keep := newMembersMap[userID]; !keep {
			if err := tx.Delete(&MemberSchema{},
				"chat_id = ? AND user_id = ?",
				chatID, userID).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
