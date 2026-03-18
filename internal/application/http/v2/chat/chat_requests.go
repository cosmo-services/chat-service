package chat_http

type SendDirectMessageRequest struct {
	RecipienUsername string `json:"recipient_username"`
	Content          string `json:"content"`
}
