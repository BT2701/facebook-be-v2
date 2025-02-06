package inbound

import (
	"net/http"
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"github.com/labstack/echo/v4"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

type PlayerHandler struct {
	playerService service.PlayerService
}

func NewPlayerHandler(playerService service.PlayerService) *PlayerHandler {
	return &PlayerHandler{playerService: playerService}
}

func (handler *PlayerHandler) CreatePlayer(c echo.Context) error {
	var player *models.Player
	player = &models.Player{} // Khởi tạo con trỏ trước khi gán giá trị
	// player.ID = primitive.NewObjectID()

	if err := c.Bind(player); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	createdPlayer, err := handler.playerService.CreatePlayer(player)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Player created successfully",
		"player":  createdPlayer,
	}, nil))
}

func (handler *PlayerHandler) GetPlayerByID(c echo.Context) error {
	playerID := c.Param("id")

	player, err := handler.playerService.GetPlayerByID(playerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"player": player,
	}, nil))
}


func (handler *PlayerHandler) UpdatePlayer(c echo.Context) error {
	// playerID := c.Param("id")
	var player *models.Player
	player = &models.Player{} // Khởi tạo con trỏ trước khi gán giá trị

	if err := c.Bind(player); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	updatedPlayer, err := handler.playerService.UpdatePlayer(player)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Player updated successfully",
		"player":  updatedPlayer,
	}, nil))
}

func (handler *PlayerHandler) DeletePlayer(c echo.Context) error {
	playerID := c.Param("id")

	err := handler.playerService.DeletePlayer(playerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Player deleted successfully",
	}, nil))
}

func (handler *PlayerHandler) GetAllPlayers(c echo.Context) error {
	players, err := handler.playerService.GetAllPlayers()
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"players": players,
	}, nil))
}

func (handler *PlayerHandler) UpdateBalance(c echo.Context) error {
	playerID := c.Param("id")
	amount := c.QueryParam("amount")

	totalAmount, err := strconv.ParseFloat(amount, 64)
	updatedBalance, err := handler.playerService.UpdateBalance(playerID, totalAmount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Balance updated successfully",
		"balance": updatedBalance,
	}, nil))
}