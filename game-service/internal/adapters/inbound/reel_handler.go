package inbound

import (
	"fmt"
	"game-service/internal/app/service"
	"game-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReelHandler struct {
	reelService *service.ReelService
}

func NewReelHandler(reelService *service.ReelService) *ReelHandler {
	return &ReelHandler{reelService: reelService}
}

func (handler *ReelHandler) GetReel(c echo.Context) error {
	reel, err := handler.reelService.LoadReel()
	if err != nil {
		response := utils.NewAPIResponse(http.StatusInternalServerError, nil, fmt.Sprintf("failed to load reel data: %v", err))
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"reel": reel,
	}, nil)
	return c.JSON(http.StatusOK, response)
}
