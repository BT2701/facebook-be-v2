package inbound

import (
	"fmt"
	"game-service/internal/app/service"
	"game-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SymbolHandler struct {
	symbolService *service.SymbolService
}

func NewSymbolHandler(symbolService *service.SymbolService) *SymbolHandler {
	return &SymbolHandler{symbolService: symbolService}
}

func (handler *SymbolHandler) GetSymbols(c echo.Context) error {
	symbols, err := handler.symbolService.LoadSymbols()
	if err != nil {
		response := utils.NewAPIResponse(http.StatusInternalServerError, nil, fmt.Sprintf("failed to load symbols: %v", err))
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := utils.NewAPIResponse(http.StatusOK, symbols,
		nil)
	return c.JSON(http.StatusOK, response)
}
