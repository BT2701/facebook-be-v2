package inbound

import (
	"fmt"
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameSessionHandler struct {
	service service.GameSessionService
}

func NewGameSessionHandler(service service.GameSessionService) *GameSessionHandler {
	return &GameSessionHandler{
		service: service,
	}
}

func (h *GameSessionHandler) CreateGameSession(c echo.Context) error {
	var session models.GameSession
	session.ID = primitive.NewObjectID()
	if err := c.Bind(&session); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}
	createdSession, err := h.service.CreateGameSession(&session)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message":     "GameSession created successfully",
		"gameSession": createdSession,
	}, nil))
}

func (h *GameSessionHandler) GetGameSessionByID(c echo.Context) error {
	id := c.Param("id")
	session, err := h.service.GetGameSessionByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"gameSession": session,
	}, nil))
}

func (h *GameSessionHandler) GetGameSessionsByPlayerID(c echo.Context) error {
	playerID := c.Param("player_id")
	page := 1
	limit := 10
	if p := c.QueryParam("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	sessions, err := h.service.GetGameSessionsByPlayerID(playerID, page, limit)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"gameSessions": sessions,
	}, nil))
}

func (h *GameSessionHandler) UpdateGameSession(c echo.Context) error {
	var session models.GameSession
	if err := c.Bind(&session); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}
	updatedSession, err := h.service.UpdateGameSession(&session)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message":     "GameSession updated successfully",
		"gameSession": updatedSession,
	}, nil))
}

func (h *GameSessionHandler) DeleteGameSession(c echo.Context) error {
	id := c.Param("id")
	err := h.service.DeleteGameSession(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "GameSession deleted successfully",
	}, nil))
}
