package chat

import "errors"

var (
	ErrChatAlreadyExsists  = errors.New("CHAT_ALREADY_EXISTS")
	ErrChatNotFound        = errors.New("CHAT_NOT_FOUND")
	ErrChatAccessDenied    = errors.New("CHAT_ACCESS_DENIED")
	ErrShortChatName       = errors.New("CHAT_NAME_SHORT")
	ErrLongChatName        = errors.New("CHAT_NAME_LONG")
	ErrManyChatMembers     = errors.New("CHAT_MANY_MEMBERS")
	ErrMessageNotFound     = errors.New("MESSAGE_NOT_FOUND")
	ErrMessageAlreadyBound = errors.New("MESSAGE_ALREADY_BOUND")
	ErrMemberNotFound      = errors.New("MEMBER_NOT_FOUND")
)
