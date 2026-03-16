package chat

type NewMessagePayload struct {
	Message   *Message
	Reply     *Message
	Chat      *Chat
	Sender    *User
	Recipient *User
}
