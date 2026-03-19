package chat

import "errors"

var (
	ErrCompanionNotPresented = errors.New("companion is not presented")
)

type ContextType string

const (
	ContextTypeMain        ContextType = "main"
	ContextTypeReplyParent ContextType = "reply_parent"
)

type ViewService struct {
}

func NewViewService() *ViewService {
	return &ViewService{}
}

func (s *ViewService) ProjectCollectionForViewers(viewersId []string, collections *Collection) (map[string]*CollectionView, error) {
	collectionViewMap := make(map[string]*CollectionView)
	for _, viewerId := range viewersId {
		collView, err := s.ProjectCollectionForViewer(viewerId, collections)
		if err != nil {
			return nil, err
		}

		collectionViewMap[viewerId] = collView
	}
	return collectionViewMap, nil
}

func (s *ViewService) ProjectCollectionForViewer(viewerId string, collection *Collection) (*CollectionView, error) {
	collView := NewCollectionView()
	collView.HasNext = collection.HasNext
	collView.HasPrev = collection.HasPrev

	for _, user := range collection.Users {
		collView.Users[user.ID] = s.projectUser(user)
	}

	for _, chat := range collection.Chats {
		chatView, err := s.projectChat(viewerId, chat, collection.Users)
		if err != nil {
			continue
		}
		collView.Chats[chat.ID] = chatView
	}

	for _, msg := range collection.Replies {
		msgView := s.projectMessage(msg, ContextTypeReplyParent)
		collView.Messages[msg.ID] = msgView
	}

	for _, msg := range collection.Messages {
		msgView := s.projectMessage(msg, ContextTypeMain)
		collView.Messages[msg.ID] = msgView
	}

	return collView, nil
}

func (s *ViewService) projectChat(userId string, chat *Chat, users map[string]*User) (*ChatView, error) {
	chatName := chat.Name
	chatAvatar := chat.AvatarUrl

	if chat.Type == ChatTypeDirect {
		companionId := chat.GetMemberIdsExcluding(userId)[0]
		companion, exists := users[companionId]

		if !exists {
			return nil, ErrCompanionNotPresented
		}

		chatName = companion.DisplayName
		chatAvatar = companion.AvatarUrl
	}

	memberViews := make([]*MemberView, len(chat.Members))
	for i, member := range chat.Members {
		memberViews[i] = &MemberView{
			UserID:   member.UserID,
			Role:     string(member.Role),
			JoinedAt: member.JoinedAt,
		}
	}

	return &ChatView{
		ID:          chat.ID,
		Type:        string(chat.Type),
		Name:        chatName,
		AvatarUrl:   chatAvatar,
		Description: chat.Description,
		CreatedAt:   chat.CreatedAt,
		Members:     memberViews,
	}, nil
}

func (s *ViewService) projectMessage(msg *Message, contextType ContextType) *MessageView {
	return &MessageView{
		ID:          msg.ID,
		ChatID:      msg.ChatID,
		SenderID:    msg.SenderID,
		Content:     msg.Content,
		ReplyToId:   msg.ReplyToId,
		Edited:      !msg.CreatedAt.Equal(msg.UpdatedAt),
		ContextType: string(contextType),
		CreatedAt:   msg.CreatedAt,
	}
}

func (s *ViewService) projectUser(user *User) *UserView {
	return &UserView{
		UserID:      user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		AvatarUrl:   user.AvatarUrl,
	}
}
