package inbound

import (
	"net/http"

	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"

	"github.com/labstack/echo/v4"
)

type SpinHandler struct {
	spin    service.SpinService
	reels   service.ReelsService
	symbols service.SymbolsService
	paylines service.PaylinesService
}

func NewSpinHandler(
	spin service.SpinService,
	reels service.ReelsService,
	symbols service.SymbolsService,
	paylines service.PaylinesService,
) *SpinHandler {
	return &SpinHandler{spin: spin, reels: reels, symbols: symbols, paylines: paylines}
}

func unwrap(common models.Common) interface{} {
	if common.Data != nil {
		return common.Data
	}
	return common
}

func (h *SpinHandler) GetDefaultReels(c echo.Context) error {
	reels, err := h.reels.GetReels("")
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"reel": unwrap(reels),
	}, nil))
}

func (h *SpinHandler) GetDefaultSymbols(c echo.Context) error {
	symbols, err := h.symbols.GetSymbol("")
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"symbols": unwrap(symbols),
	}, nil))
}

func (h *SpinHandler) GetDefaultPaylines(c echo.Context) error {
	paylines, err := h.paylines.GetPayline("")
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"paylines": unwrap(paylines),
	}, nil))
}

func (h *SpinHandler) GetBets(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"bet_options": h.spin.BetOptions(),
	}, nil))
}

func (h *SpinHandler) CalculateWinnings(c echo.Context) error {
	var input struct {
		FinalSymbols [][]string `json:"finalSymbols"`
		BetAmount    float64    `json:"betAmount"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	winAmount, results, err := h.spin.CalculateWinnings(input.FinalSymbols, input.BetAmount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"winAmount": winAmount,
		"results":   results,
	}, nil))
}
