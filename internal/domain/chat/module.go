package chat

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserProfileSync),
	fx.Provide(NewChatFactory),
	fx.Provide(NewMessageFactory),
	fx.Provide(NewChatService),
)
