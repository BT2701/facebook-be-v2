package inbound

import (
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PaylinesHandler struct {
	service service.PaylinesService
}

func NewPaylinesHandler(service service.PaylinesService) *PaylinesHandler {
	return &PaylinesHandler{
		service: service,
	}
}

func (h *PaylinesHandler) GetPaylines(c echo.Context) error {
	gameName := c.Param("game_name")
	paylines, err := h.service.GetPayline(gameName)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"paylines": paylines,
	}, nil))
}

func (h *PaylinesHandler) SetPaylines(c echo.Context) error {
	gameName := c.Param("game_name")
	var value models.Common
	if err := c.Bind(&value); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}
	err := h.service.SetPayline(gameName, value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Paylines updated successfully",
	}, nil))
}
