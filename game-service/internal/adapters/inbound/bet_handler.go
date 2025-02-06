package inbound

import (
	"fmt"
	"game-service/internal/app/service"
	"game-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BetHandler struct {
	betService *service.BetService
}

func NewBetHandler(betService *service.BetService) *BetHandler {
	return &BetHandler{betService: betService}
}

func (handler *BetHandler) GetBets(c echo.Context) error {
	betOptions, err := handler.betService.LoadBetOptions()
	if err != nil {
		response := utils.NewAPIResponse(http.StatusInternalServerError, nil, fmt.Sprintf("failed to load bet options: %v", err))
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := utils.NewAPIResponse(http.StatusOK, betOptions,
		nil)
	return c.JSON(http.StatusOK, response)
}
