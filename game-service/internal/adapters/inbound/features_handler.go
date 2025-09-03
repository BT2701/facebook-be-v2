package inbound

import (
	"game-service/internal/app/service"
	"game-service/internal/models"
	"game-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FeaturesHandler struct {
	service service.FeaturesService
}

func NewFeaturesHandler(service service.FeaturesService) *FeaturesHandler {
	return &FeaturesHandler{
		service: service,
	}
}

func (h *FeaturesHandler) GetFeature(c echo.Context) error {
	gameName := c.Param("game_name")
	feature, err := h.service.GetFeature(gameName)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"feature": feature,
	}, nil))
}

func (h *FeaturesHandler) SetFeature(c echo.Context) error {
	gameName := c.Param("game_name")
	var value models.Common
	if err := c.Bind(&value); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}
	err := h.service.SetFeature(gameName, value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Feature updated successfully",
	}, nil))
}
