package chat_http

import (
	_ "main/docs"
	"main/internal/application/http/v2/auth"

	"main/pkg"
)

type ChatRoutes struct {
	handler        pkg.RequestHandler
	controller     *ChatController
	authMiddleware *auth.AuthMiddleware
}

func NewChatRoutes(handler pkg.RequestHandler, controller *ChatController, authMiddleware *auth.AuthMiddleware) *ChatRoutes {
	return &ChatRoutes{
		handler:        handler,
		controller:     controller,
		authMiddleware: authMiddleware,
	}
}

func (r *ChatRoutes) Setup() {
	api := r.handler.Gin.Group("/api/v2/chat")

	read := api.Group("/")
	read.Use(r.authMiddleware.RequireAuth())
	{
		read.GET("/history", r.controller.GetChatsHistory)
		read.GET("/direct/:username/messages", r.controller.GetDirectMessages)
		read.GET("/:chat_id/messages", r.controller.GetChatMessages)
	}

	cmd := api.Group("/")
	cmd.Use(r.authMiddleware.RequireAuth(), r.authMiddleware.RequireActive())
	{
		cmd.POST("/direct/:username/messages", r.controller.SendDirectMessage)
	}
}
