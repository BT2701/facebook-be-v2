package inbound

import (
	"net/http"
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BonusGameHandler struct {
	bonusGameService service.BonusGameService
}

func NewBonusGameHandler(bonusGameService service.BonusGameService) *BonusGameHandler {
	return &BonusGameHandler{bonusGameService: bonusGameService}
}

func (handler *BonusGameHandler) CreateBonusGame(c echo.Context) error {
	var bonusGame *models.BonusGame
	bonusGame = &models.BonusGame{} // Khởi tạo con trỏ trước khi gán giá trị
	bonusGame.ID = primitive.NewObjectID()

	if err := c.Bind(bonusGame); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	createdBonusGame, err := handler.bonusGameService.CreateBonusGame(bonusGame)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "BonusGame created successfully",
		"bonusGame":  createdBonusGame,
	}, nil))
}

func (handler *BonusGameHandler) GetBonusGameByID(c echo.Context) error {
	bonusGameID := c.Param("id")

	bonusGame, err := handler.bonusGameService.GetBonusGameByID(bonusGameID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"bonusGame": bonusGame,
	}, nil))
}

func (handler *BonusGameHandler) GetBonusGamesByPlayerID(c echo.Context) error {
	playerID := c.Param("player_id")

	bonusGames, err := handler.bonusGameService.GetBonusGamesByPlayerID(playerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"bonusGames": bonusGames,
	}, nil))
}

func (handler *BonusGameHandler) UpdateBonusGame(c echo.Context) error {
	// bonusGameID := c.Param("id")
	return nil
}

func (handler *BonusGameHandler) DeleteBonusGame(c echo.Context) error {
	bonusGameID := c.Param("id")

	err := handler.bonusGameService.DeleteBonusGame(bonusGameID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "BonusGame deleted successfully",
	}, nil))
}
