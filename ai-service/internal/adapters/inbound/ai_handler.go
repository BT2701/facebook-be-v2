package inbound

import (
	"ai-service/internal/app/service"
	"ai-service/pkg/utils"
	"net/http"
	"strings"

	"github.com/BT2701/facebook-be-v2/shared/auth"
	"github.com/labstack/echo/v4"
)

type AIHandler struct {
	service service.AIService
}

func NewAIHandler(aiService service.AIService) *AIHandler {
	return &AIHandler{service: aiService}
}

type chatRequest struct {
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
}

func (handler *AIHandler) GetInbox(c echo.Context) error {
	userID := resolveUserID(c, "")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "user_id is required"))
	}

	conversation, messages, err := handler.service.GetInbox(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"conversation": conversation,
		"messages":     messages,
	}, nil))
}

func (handler *AIHandler) Chat(c echo.Context) error {
	var input chatRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	userID := resolveUserID(c, input.UserID)
	conversation, userMessage, assistantMessage, err := handler.service.Chat(
		c.Request().Context(),
		userID,
		input.ConversationID,
		input.Message,
	)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "too long") {
			status = http.StatusBadRequest
		}
		return c.JSON(status, utils.NewAPIResponse(status, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"conversation": conversation,
		"user_message": userMessage,
		"assistant":    assistantMessage,
	}, nil))
}

func (handler *AIHandler) Clear(c echo.Context) error {
	userID := resolveUserID(c, "")
	if err := handler.service.Clear(c.Request().Context(), userID, c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Conversation cleared",
	}, nil))
}

func resolveUserID(c echo.Context, fallback string) string {
	if id := auth.UserID(c); id != "" {
		return id
	}
	return strings.TrimSpace(fallback)
}
