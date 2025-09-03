package inbound

import (
	"game-service/internal/app/service"
	"game-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BackupHandler struct {
	service service.BackupService
}

func NewBackupHandler(service service.BackupService) *BackupHandler {
	return &BackupHandler{service: service}
}

func (h *BackupHandler) BackupAll(c echo.Context) error {
	err := h.service.BackupAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Backup completed successfully",
	}, nil))
}
