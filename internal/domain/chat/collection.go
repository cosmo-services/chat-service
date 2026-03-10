package chat

type Collection struct {
	Chats    map[string]*Chat    `json:"chats"`
	Messages map[string]*Message `json:"messages"`
	Replies  map[string]*Message `json:"replies"`
	Users    map[string]*User    `json:"users"`

	HasNext bool `json:"has_next"`
	HasPrev bool `json:"has_prev"`
}
