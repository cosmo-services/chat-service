package social

type SocialClient interface {
	GetProfileByUserId(userId string) (*UserProfile, error)
	GetProfileByUsername(username string) (*UserProfile, error)
}
