package chat

type Collection struct {
	Chats    map[string]*Chat    `json:"chats"`
	Messages map[string]*Message `json:"messages"`
	Replies  map[string]*Message `json:"replies"`
	Users    map[string]*User    `json:"users"`

	HasNext bool `json:"has_next"`
	HasPrev bool `json:"has_prev"`
}

func NewCollection() *Collection {
	return &Collection{
		Chats:    make(map[string]*Chat),
		Messages: make(map[string]*Message),
		Replies:  make(map[string]*Message),
		Users:    make(map[string]*User),
	}
}

func NewDirectMessageCollection(direct *Chat, msg *Message, sender *User, recipient *User) *Collection {
	return &Collection{
		Chats:    map[string]*Chat{direct.ID: direct},
		Messages: map[string]*Message{msg.ID: msg},
		Users: map[string]*User{
			sender.ID:    sender,
			recipient.ID: recipient,
		},
		Replies: make(map[string]*Message),
	}
}
