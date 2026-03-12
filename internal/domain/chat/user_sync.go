package chat

import "main/internal/domain/social"

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

func (s *UserProfileSync) EnsureUserExistsById(userId string) error {
	exists, err := s.userRepo.UserExistsById(userId)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	userProfile, err := s.socialClient.GetProfileByUserId(userId)
	if err != nil {
		return err
	}

	user := MapProfileToUser(userProfile)
	if err := s.userRepo.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *UserProfileSync) EnsureUserExistsByUsername(username string) error {
	exists, err := s.userRepo.UserExistsByUsername(username)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	userProfile, err := s.socialClient.GetProgileByUsername(username)
	if err != nil {
		return err
	}

	user := MapProfileToUser(userProfile)
	if err := s.userRepo.CreateUser(user); err != nil {
		return err
	}

	return nil
}
