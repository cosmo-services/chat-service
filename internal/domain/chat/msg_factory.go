package chat

type MessageFactory struct {
}

func NewMessageFactory() *MessageFactory {
	return &MessageFactory{}
}

func (f *MessageFactory) CreateTextMessage(senderId, content string) (*Message, error) {
	return &Message{
		SenderID: senderId,
		Content:  content,
	}, nil
}
