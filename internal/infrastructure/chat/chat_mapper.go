package chat_infrastructure

import (
	chat_domain "main/internal/domain/chat"
)

func ToDomainChat(s *ChatSchema) *chat_domain.Chat {
	if s == nil {
		return nil
	}

	members := make([]*chat_domain.ChatMember, 0, len(s.Members))
	for _, m := range s.Members {
		if m != nil {
			members = append(members, ToDomainMember(m))
		}
	}

	return &chat_domain.Chat{
		ID:          s.ID,
		Type:        chat_domain.ChatType(s.Type),
		Name:        s.Name,
		Description: s.Description,
		CreatedBy:   s.CreatedBy,
		Members:     members,
		Deleted:     !s.DeletedAt.Time.IsZero(),
		AvatarUrl:   s.AvatarUrl,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		DeletedAt:   s.DeletedAt.Time,
	}
}

func ToSchemaChat(c *chat_domain.Chat) *ChatSchema {
	if c == nil {
		return nil
	}

	members := make([]*MemberSchema, 0, len(c.Members))
	for _, m := range c.Members {
		if m != nil {
			members = append(members, ToSchemaMember(m))
		}
	}

	return &ChatSchema{
		ID:          c.ID,
		Type:        string(c.Type),
		Name:        c.Name,
		Description: c.Description,
		CreatedBy:   c.CreatedBy,
		Members:     members,
		AvatarUrl:   c.AvatarUrl,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func ToDomainMember(s *MemberSchema) *chat_domain.ChatMember {
	if s == nil {
		return nil
	}

	return &chat_domain.ChatMember{
		UserID:    s.UserID,
		Role:      chat_domain.ChatMemberRole(s.Role),
		JoinedAt:  s.JoinedAt,
		DeletedAt: s.DeletedAt.Time,
	}
}

func ToSchemaMember(m *chat_domain.ChatMember) *MemberSchema {
	if m == nil {
		return nil
	}

	return &MemberSchema{
		UserID:   m.UserID,
		Role:     string(m.Role),
		JoinedAt: m.JoinedAt,
	}
}
