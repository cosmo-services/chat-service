package chat

type MessageSearchFilter struct {
	Cursor    string
	Direction string
	Limit     int
}

type ChatSearchFilter struct {
	Cursor    string
	Direction string
	Limit     int
}
