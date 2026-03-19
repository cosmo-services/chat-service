package chat

import "time"

type ChatFactory struct {
}

func NewChatFactory() *ChatFactory {
	return &ChatFactory{}
}

func (f *ChatFactory) CreateDirectChat(user1ID, user2ID string) (*Chat, error) {
	return &Chat{
		Type: ChatTypeDirect,
		Members: []*ChatMember{
			{UserID: user1ID, Role: RoleMember},
			{UserID: user2ID, Role: RoleMember},
		},
		CreatedAt: time.Now(),
	}, nil
}

func (f *ChatFactory) CreateGroupChat(creatorID, name string, initialMembers []string) (*Chat, error) {
	members := make([]*ChatMember, 0, len(initialMembers)+1)

	members = append(members, &ChatMember{
		UserID: creatorID,
		Role:   RoleAdmin,
	})

	for _, userID := range initialMembers {
		if userID == creatorID {
			continue
		}

		members = append(members, &ChatMember{
			UserID: userID,
			Role:   RoleMember,
		})
	}

	return &Chat{
		Type:      ChatTypeGroup,
		Name:      name,
		Members:   members,
		CreatedBy: creatorID,
	}, nil
}
