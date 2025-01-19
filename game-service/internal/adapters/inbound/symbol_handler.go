package inbound

import (
	"net/http"
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SymbolHandler struct {
	symbolService service.SymbolService
}

func NewSymbolHandler(symbolService service.SymbolService) *SymbolHandler {
	return &SymbolHandler{symbolService: symbolService}
}

func (handler *SymbolHandler) CreateSymbol(c echo.Context) error {
	var symbol *models.Symbol
	symbol = &models.Symbol{} // Khởi tạo con trỏ trước khi gán giá trị
	symbol.ID = primitive.NewObjectID()

	if err := c.Bind(symbol); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	createdSymbol, err := handler.symbolService.CreateSymbol(symbol)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Symbol created successfully",
		"symbol":  createdSymbol,
	}, nil))
}

func (handler *SymbolHandler) GetSymbolByID(c echo.Context) error {
	symbolID := c.Param("id")

	symbol, err := handler.symbolService.GetSymbolByID(symbolID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"symbol": symbol,
	}, nil))
}

func (handler *SymbolHandler) GetSymbolsByGameID(c echo.Context) error {
	gameID := c.Param("game_id")

	symbols, err := handler.symbolService.GetSymbolsByGameID(gameID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"symbols": symbols,
	}, nil))
}

