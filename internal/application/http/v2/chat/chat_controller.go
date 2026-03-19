package chat_http

import (
	"errors"
	chat_domain "main/internal/domain/chat"
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chatService *chat_domain.ChatService
	logger      pkg.Logger
}

func NewChatController(
	chatService *chat_domain.ChatService,
	logger pkg.Logger,
) *ChatController {
	return &ChatController{
		chatService: chatService,
		logger:      logger,
	}
}

// GetChatsHistory godoc
//
// @Summary Get user chats
// @Description Get cursor-paginated list of user's chats sorted by last activity (updated_at)
// @Tags chat
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param cursor query string false "Cursor (timestamp, e.g., 2024-01-15T10:30:00Z). First page if empty"
// @Param direction query string false "Pagination direction: 'next' (newer chats) or 'prev' (older chats). Default is 'prev' (most recent first)."
// @Param limit query int false "Number of chats per page (default: 10, max: 100)"
// @Success 200 {object} chat.CollectionView
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Failure 404 {object} map[string]string "User not found"
// @Router /history [get]
func (controller *ChatController) GetChatsHistory(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	cursor := ctx.Query("cursor")
	direction := ctx.DefaultQuery("direction", "prev")

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	collection, err := controller.chatService.GetChatsHistoryView(userId, chat_domain.ChatPageFilter{
		Cursor:    cursor,
		Direction: direction,
		Limit:     limit,
	})
	if err != nil {
		if errors.Is(err, chat_domain.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, collection)
}

// GetDirectMessages godoc
//
// @Summary Get direct message history
// @Description Get cursor-paginated messages from a direct chat with specific user
// @Tags chat
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Recipient username (the user you're chatting with)"
// @Param cursor query string false "Cursor - last message ID from previous page (e.g., 'msg_1234567890'). Empty for first page."
// @Param direction query string false "Pagination direction: 'next' (older messages) or 'prev' (newer messages). Default is 'prev' (most recent first)."
// @Param limit query int false "Number of messages per page (default: 20, max: 100)"
// @Success 200 {object} chat.CollectionView
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Failure 404 {object} map[string]string "User not found"
// @Router /direct/{username}/messages [get]
func (controller *ChatController) GetDirectMessages(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	cursor := ctx.Query("cursor")
	direction := ctx.DefaultQuery("direction", "prev")
	username := ctx.Param("username")

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	collection, err := controller.chatService.GetDirectMessageHistoryView(userId, username,
		chat_domain.MessagePageFilter{
			Cursor:    cursor,
			Direction: direction,
			Limit:     limit,
		})
	if err != nil {
		if errors.Is(err, chat_domain.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, collection)
}

// GetChatMessages godoc
//
// @Summary Get chat message history
// @Description Get cursor-paginated messages from a specific chat (group or direct)
// @Tags chat
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param chat_id path string true "Chat ID (e.g., 'chat_1234567890')"
// @Param cursor query string false "Cursor - last message ID from previous page (e.g., 'msg_1234567890'). Empty for first page."
// @Param direction query string false "Pagination direction: 'next' (older messages) or 'prev' (newer messages). Default is 'prev' (most recent first)."
// @Param limit query int false "Number of messages per page (default: 20, max: 100)"
// @Success 200 {object} chat.CollectionView
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Failure 404 {object} map[string]string "Chat not found"
// @Router /{chat_id}/messages [get]
func (controller *ChatController) GetChatMessages(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	cursor := ctx.Query("cursor")
	direction := ctx.DefaultQuery("direction", "prev")
	chatId := ctx.Param("chat_id")

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	collection, err := controller.chatService.GetChatMessageHistoryView(userId, chatId,
		chat_domain.MessagePageFilter{
			Cursor:    cursor,
			Direction: direction,
			Limit:     limit,
		})
	if err != nil {
		if errors.Is(err, chat_domain.ErrChatNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, collection)
}

// SendDirectMessage godoc
//
// @Summary Send direct message
// @Description Send a message to a user by username. Creates or updates direct chat automatically
// @Tags chat
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Recipient username (the user you want to message)"
// @Param request body SendMessageRequest true "Message content (text, media, etc.)"
// @Success 201 {object} map[string]string "Message sent successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid message format"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found - recipient doesn't exist"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /direct/{username}/messages [post]
func (controller *ChatController) SendDirectMessage(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	username := ctx.Param("username")

	var req *SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		controller.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := controller.chatService.SendDirectMessage(userId, username, req.Content); err != nil {
		if errors.Is(err, chat_domain.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "message sent",
	})
}

// SendChatMessage godoc
//
// @Summary Send message to chat
// @Description Send a message to an existing group chat by chat ID
// @Tags chat
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param chat_id path string true "Chat ID (e.g., 'chat_1234567890')"
// @Param request body SendMessageRequest true "Message content (text, reply_to, attachments)"
// @Success 201 {object} map[string]string "Message sent successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid message format or empty content"
// @Failure 401 {object} map[string]string "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]string "Forbidden - user is not a member of this chat"
// @Failure 404 {object} map[string]string "Chat not found - chat doesn't exist or has been deleted"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /{chat_id}/messages [post]
func (controller *ChatController) SendChatMessage(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	chatId := ctx.Param("chat_id")

	var req *SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		controller.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := controller.chatService.SendChatMessage(userId, chatId, req.Content); err != nil {
		if errors.Is(err, chat_domain.ErrChatNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "message sent",
	})
}
