package inbound

import (
	"net/http"
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BetHandler struct {
	betService service.BetService
}

func NewBetHandler(betService service.BetService) *BetHandler {
	return &BetHandler{betService: betService}
}

func (handler *BetHandler) CreateBet(c echo.Context) error {
	var bet *models.BetOption
	bet = &models.BetOption{} // Khởi tạo con trỏ trước khi gán giá trị
	bet.ID = primitive.NewObjectID()

	if err := c.Bind(bet); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	createdBet, err := handler.betService.CreateBet(bet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Bet created successfully",
		"bet":  createdBet,
	}, nil))
}

func (handler *BetHandler) GetBetByID(c echo.Context) error {
	betID := c.Param("id")

	bet, err := handler.betService.GetBetByID(betID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"bet": bet,
	}, nil))
}

func (handler *BetHandler) GetBetsByPlayerID(c echo.Context) error {
	playerID := c.Param("player_id")

	bets, err := handler.betService.GetBetsByPlayerID(playerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"bets": bets,
	}, nil))
}

func (handler *BetHandler) UpdateBet(c echo.Context) error {
	// betID := c.Param("id")
	var bet *models.BetOption
	bet = &models.BetOption{} // Khởi tạo con trỏ trước khi gán giá trị

	if err := c.Bind(bet); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	updatedBet, err := handler.betService.UpdateBet(bet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Bet updated successfully",
		"bet": updatedBet,
	}, nil))
}

func (handler *BetHandler) DeleteBet(c echo.Context) error {
	betID := c.Param("id")

	err := handler.betService.DeleteBet(betID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Bet deleted successfully",
	}, nil))
}

