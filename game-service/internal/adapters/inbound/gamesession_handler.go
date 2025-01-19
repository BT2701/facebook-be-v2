package inbound

import (
	"net/http"
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameSessionHandler struct {
	gameSessionService service.GameSessionService
}

func NewGameSessionHandler(gameSessionService service.GameSessionService) *GameSessionHandler {
	return &GameSessionHandler{gameSessionService: gameSessionService}
}

func (handler *GameSessionHandler) CreateGameSession(c echo.Context) error {
	var gameSession *models.GameSession
	gameSession = &models.GameSession{} // Khởi tạo con trỏ trước khi gán giá trị
	gameSession.ID = primitive.NewObjectID()

	if err := c.Bind(gameSession); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	createdGameSession, err := handler.gameSessionService.CreateGameSession(gameSession)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "GameSession created successfully",
		"gameSession":  createdGameSession,
	}, nil))
}

func (handler *GameSessionHandler) GetGameSessionByID(c echo.Context) error {
	gameSessionID := c.Param("id")

	gameSession, err := handler.gameSessionService.GetGameSessionByID(gameSessionID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"gameSession": gameSession,
	}, nil))
}

func (handler *GameSessionHandler) GetGameSessionsByPlayerID(c echo.Context) error {
	playerID := c.Param("player_id")

	gameSessions, err := handler.gameSessionService.GetGameSessionsByPlayerID(playerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"gameSessions": gameSessions,
	}, nil))
}

func (handler *GameSessionHandler) UpdateGameSession(c echo.Context) error {
	// gameSessionID := c.Param("id")
	return nil
}

func (handler *GameSessionHandler) DeleteGameSession(c echo.Context) error {
	gameSessionID := c.Param("id")

	err := handler.gameSessionService.DeleteGameSession(gameSessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "GameSession deleted successfully",
	}, nil))
}
