package inbound

import (
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ConfigsHandler struct {
	service service.ConfigsService
}

func NewConfigsHandler(service service.ConfigsService) *ConfigsHandler {
	return &ConfigsHandler{
		service: service,
	}
}

func (h *ConfigsHandler) GetConfig(c echo.Context) error {
	gameName := c.Param("game_name")
	config, err := h.service.GetConfig(gameName)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"config": config,
	}, nil))
}

func (h *ConfigsHandler) SetConfig(c echo.Context) error {
	gameName := c.Param("game_name")
	var value models.Common
	if err := c.Bind(&value); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}
	err := h.service.SetConfig(gameName, value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Config updated successfully",
	}, nil))
}
