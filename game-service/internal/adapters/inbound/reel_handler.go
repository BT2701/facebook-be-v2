package inbound

import (
	"net/http"
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReelHandler struct {
	reelService service.ReelService
}

func NewReelHandler(reelService service.ReelService) *ReelHandler {
	return &ReelHandler{reelService: reelService}
}

func (handler *ReelHandler) CreateReel(c echo.Context) error {
	var reel *models.Reel
	reel = &models.Reel{} // Khởi tạo con trỏ trước khi gán giá trị
	reel.ID = primitive.NewObjectID()

	if err := c.Bind(reel); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	createdReel, err := handler.reelService.CreateReel(reel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Reel created successfully",
		"reel":  createdReel,
	}, nil))
}

func (handler *ReelHandler) GetReelByID(c echo.Context) error {
	reelID := c.Param("id")

	reel, err := handler.reelService.GetReelByID(reelID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"reel": reel,
	}, nil))
}

func (handler *ReelHandler) GetReelsByGameID(c echo.Context) error {
	gameID := c.Param("game_id")

	reels, err := handler.reelService.GetReelsByGameID(gameID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"reels": reels,
	}, nil))
}

func (handler *ReelHandler) GetReelsByGameSessionID(c echo.Context) error {
	// gameSessionID := c.Param("game_session_id")

	// reels, err := handler.reelService.GetReelsByGameSessionID(gameSessionID)
	// if err != nil {
	// 	return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	// }

	// return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
	// 	"reels": reels,
	// }, nil))
	return nil
}

func (handler *ReelHandler) UpdateReel(c echo.Context) error {
	// reelID := c.Param("id")
	return nil
}

func (handler *ReelHandler) DeleteReel(c echo.Context) error {
	reelID := c.Param("id")

	err := handler.reelService.DeleteReel(reelID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Reel deleted successfully",
	}, nil))
}

