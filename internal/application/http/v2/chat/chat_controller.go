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
// @Description Get cursor paginated chats
// @Tags chat
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param cursor query string false "Cursor"
// @Param direction query string false "Direction"
// @Param limit query string false "limit"
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
// @Summary Get direct chat history
// @Description Get cursor paginated messages
// @Tags chat
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param cursor query string false "Cursor"
// @Param direction query string false "Direction"
// @Param limit query string false "limit"
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
