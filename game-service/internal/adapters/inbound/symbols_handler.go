package inbound

import (
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SymbolsHandler struct {
	service service.SymbolsService
}

func NewSymbolsHandler(service service.SymbolsService) *SymbolsHandler {
	return &SymbolsHandler{
		service: service,
	}
}

func (h *SymbolsHandler) GetSymbols(c echo.Context) error {
	gameName := c.Param("game_name")
	symbols, err := h.service.GetSymbol(gameName)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"symbols": symbols,
	}, nil))
}

func (h *SymbolsHandler) SetSymbols(c echo.Context) error {
	gameName := c.Param("game_name")
	var value models.Common
	if err := c.Bind(&value); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}
	err := h.service.SetSymbol(gameName, value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Symbols updated successfully",
	}, nil))
}
