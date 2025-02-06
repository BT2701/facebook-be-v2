package inbound

import (
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
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
	gameSession.CreatedAt = time.Now()

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

	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 7 // Mặc định 7 dòng mỗi trang
	}

	gameSessions, err := handler.gameSessionService.GetGameSessionsByPlayerID(playerID, page, limit)
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
