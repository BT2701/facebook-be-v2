package inbound

import (
	"fmt"
	"game-service/internal/app/service"
	"game-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PaylineHandler struct {
	paylineService *service.PaylineService
}

func NewPaylineHandler(paylineService *service.PaylineService) *PaylineHandler {
	return &PaylineHandler{paylineService: paylineService}
}

func (handler *PaylineHandler) GetPaylines(c echo.Context) error {
	paylines, err := handler.paylineService.LoadPaylines()
	if err != nil {
		response := utils.NewAPIResponse(http.StatusInternalServerError, nil, fmt.Sprintf("failed to load paylines: %v", err))
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := utils.NewAPIResponse(http.StatusOK,
		paylines,
		nil)
	return c.JSON(http.StatusOK, response)
}

func (handler *PaylineHandler) CalculateWinnings(c echo.Context) error {
	var request struct {
		FinalSymbols [][]string `json:"finalSymbols"`
		BetAmount    float64    `json:"betAmount"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid request body"))
	}

	winAmount, results := handler.paylineService.CalculateWinnings(request.FinalSymbols, request.BetAmount)
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"winAmount": winAmount,
		"results":   results,
	}, nil))
}
