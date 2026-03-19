package chat

import (
	"errors"
	"main/internal/domain/social"
)

type UserProfileSync struct {
	userRepo     UserRepository
	socialClient social.SocialClient
}

func NewUserProfileSync(
	userRepo UserRepository,
	socialClient social.SocialClient,
) *UserProfileSync {
	return &UserProfileSync{
		userRepo:     userRepo,
		socialClient: socialClient,
	}
}

func (s *UserProfileSync) EnsureUserExists(opts UserSearchOptions) error {
	var userProfile *social.UserProfile
	var exists bool
	var err error

	if opts.UserID != "" {
		exists, err = s.userRepo.UserExistsById(opts.UserID)
	} else if opts.Username != "" {
		exists, err = s.userRepo.UserExistsByUsername(opts.Username)
	} else {
		return ErrUserNotFound
	}

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	if opts.UserID != "" {
		userProfile, err = s.socialClient.GetProfileByUserId(opts.UserID)
	} else if opts.Username != "" {
		userProfile, err = s.socialClient.GetProfileByUsername(opts.Username)
	}

	if err != nil {
		return err
	}

	user := s.mapProfileToUser(userProfile)
	if err := s.userRepo.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *UserProfileSync) GetUser(opts UserSearchOptions) (*User, error) {
	var userProfile *social.UserProfile
	var user *User
	var err error

	if opts.UserID != "" {
		user, err = s.userRepo.GetUserById(opts.UserID)
	} else if opts.Username != "" {
		user, err = s.userRepo.GetUserByUsername(opts.Username)
	} else {
		return nil, ErrUserNotFound
	}

	if err == nil {
		return user, err
	}

	if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	if opts.UserID != "" {
		userProfile, err = s.socialClient.GetProfileByUserId(opts.UserID)
	} else if opts.Username != "" {
		userProfile, err = s.socialClient.GetProfileByUsername(opts.Username)
	}

	if err != nil {
		return nil, err
	}

	syncUser := s.mapProfileToUser(userProfile)
	if err := s.userRepo.CreateUser(syncUser); err != nil {
		return nil, err
	}

	return syncUser, nil
}

func (s *UserProfileSync) mapProfileToUser(profile *social.UserProfile) *User {
	return &User{
		ID:          profile.UserID,
		Username:    profile.Username,
		DisplayName: profile.DisplayName,
		AvatarUrl:   profile.AvatarUrl,
		IsActive:    profile.IsActive,
	}
}
