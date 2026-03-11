package chat_infrastructure

import (
	"errors"

	chat_domain "main/internal/domain/chat"
	"main/pkg"

	"gorm.io/gorm"
)

type userRepository struct {
	db pkg.GormDB
}

func NewUserRepository(db pkg.GormDB) chat_domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *chat_domain.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	schema := ToSchemaUser(user)
	return r.db.DB.Create(schema).Error
}

func (r *userRepository) GetUserById(userID string) (*chat_domain.User, error) {
	var userSchema UserSchema

	err := r.db.DB.First(&userSchema, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, chat_domain.ErrUserNotFound
		}
		return nil, err
	}

	return ToDomainUser(&userSchema), nil
}

func (r *userRepository) GetUserByUsername(username string) (*chat_domain.User, error) {
	var userSchema UserSchema

	err := r.db.DB.Where("username = ?", username).First(&userSchema).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, chat_domain.ErrUserNotFound
		}
		return nil, err
	}

	return ToDomainUser(&userSchema), nil
}

func (r *userRepository) UpdateUser(user *chat_domain.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	schema := ToSchemaUser(user)

	return r.db.DB.Model(&UserSchema{}).
		Where("user_id = ?", user.UserID).
		Updates(map[string]interface{}{
			"username":     schema.Username,
			"display_name": schema.DisplayName,
			"avatar_url":   schema.AvatarUrl,
			"updated_at":   schema.UpdatedAt,
		}).Error
}

func (r *userRepository) DeleteUserByUsername(username string) error {
	result := r.db.DB.Delete(&UserSchema{}, "username = ?", username)

	if result.RowsAffected == 0 {
		return chat_domain.ErrUserNotFound
	}
	return result.Error
}

func (r *userRepository) DeleteUserById(userID string) error {
	result := r.db.DB.Delete(&UserSchema{}, "user_id = ?", userID)

	if result.RowsAffected == 0 {
		return chat_domain.ErrUserNotFound
	}
	return result.Error
}

func (r *userRepository) UserExistsById(userID string) (bool, error) {
	var count int64
	err := r.db.DB.Model(&UserSchema{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	return count > 0, err
}

func (r *userRepository) UserExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.DB.Model(&UserSchema{}).
		Where("username = ?", username).
		Count(&count).Error

	return count > 0, err
}

func (r *userRepository) GetUserByIdUnscoped(userID string) (*chat_domain.User, error) {
	var userSchema UserSchema

	err := r.db.DB.Unscoped().First(&userSchema, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, chat_domain.ErrUserNotFound
		}
		return nil, err
	}

	return ToDomainUser(&userSchema), nil
}

func (r *userRepository) RestoreUser(userID string) error {
	return r.db.DB.Unscoped().
		Model(&UserSchema{}).
		Where("user_id = ?", userID).
		Update("deleted_at", nil).Error
}
