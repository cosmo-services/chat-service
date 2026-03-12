package social

type SocialClient interface {
	GetProfileByUserId(userId string) (*UserProfile, error)
	GetProgileByUsername(username string) (*UserProfile, error)
}
