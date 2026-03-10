package chat

type User struct {
	UserId      string `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarUrl   string `json:"avatar_url"`
}
