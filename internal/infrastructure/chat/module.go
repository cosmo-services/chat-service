package chat_infrastructure

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewChatQuery),
	fx.Provide(NewChatRepository),
	fx.Provide(NewMessageRepository),
	fx.Provide(NewUserRepository),
	fx.Provide(NewPublisher),
)
