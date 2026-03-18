package chat

type MessagePageFilter struct {
	Cursor    string
	Direction string
	Limit     int
}

type ChatPageFilter struct {
	Cursor    string
	Direction string
	Limit     int
}
