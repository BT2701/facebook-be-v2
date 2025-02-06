package inbound

import (
	"net/http"
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameResultHandler struct {

	gameResultService service.GameResultService
}

func NewGameResultHandler(gameResultService service.GameResultService) *GameResultHandler {
	return &GameResultHandler{gameResultService: gameResultService}
}

func (handler *GameResultHandler) CreateGameResult(c echo.Context) error {
	var gameResult *models.GameResult
	gameResult = &models.GameResult{} // Khởi tạo con trỏ trước khi gán giá trị
	gameResult.ID = primitive.NewObjectID()

	if err := c.Bind(gameResult); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	createdGameResult, err := handler.gameResultService.CreateGameResult(gameResult)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "GameResult created successfully",
		"gameResult":  createdGameResult,
	}, nil))
}

func (handler *GameResultHandler) GetGameResultByID(c echo.Context) error {
	gameResultID := c.Param("id")

	gameResult, err := handler.gameResultService.GetGameResultByID(gameResultID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"gameResult": gameResult,
	}, nil))
}

func (handler *GameResultHandler) GetGameResultsByPlayerID(c echo.Context) error {
	playerID := c.Param("player_id")

	gameResults, err := handler.gameResultService.GetGameResultsByPlayerID(playerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"gameResults": gameResults,
	}, nil))
}

func (handler *GameResultHandler) UpdateGameResult(c echo.Context) error {
	// gameResultID := c.Param("id")
	return nil
}

func (handler *GameResultHandler) DeleteGameResult(c echo.Context) error {
	gameResultID := c.Param("id")

	err := handler.gameResultService.DeleteGameResult(gameResultID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "GameResult deleted successfully",
	}, nil))
}

func (handler *GameResultHandler) GetGameResultsBySessionID(c echo.Context) error {
	sessionID := c.Param("session_id")

	gameResults, err := handler.gameResultService.GetGameResultsBySessionID(sessionID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"gameResults": gameResults,
	}, nil))
}
