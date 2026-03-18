package chat_http

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewChatController),
	fx.Provide(NewChatRoutes),
)
