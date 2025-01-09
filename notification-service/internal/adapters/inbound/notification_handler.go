package inbound

import (
	"net/http"
	"notification-service/internal/model"
	"notification-service/pkg/utils"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"notification-service/internal/app/services"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// CreateNotification handles creating a new notification
func (handler *NotificationHandler) CreateNotification(c echo.Context) error {
	var notification model.Notification
	if err := c.Bind(&notification); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	notification.ID = primitive.NewObjectID()
	notification.Timeline = time.Now()

	createdNotification, err := handler.notificationService.CreateNotification(&notification)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusCreated, utils.NewAPIResponse(http.StatusCreated, map[string]interface{}{
		"message": "Notification created successfully",
		"notification": createdNotification,
	}, nil))
}

// GetNotification handles retrieving a notification by ID
func (handler *NotificationHandler) GetNotification(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Missing notification ID"))
	}

	notification, err := handler.notificationService.GetNotification(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	if notification == nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, "Notification not found"))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"notification": notification,
	}, nil))
}

// UpdateNotification handles updating a notification by ID
func (handler *NotificationHandler) UpdateNotification(c echo.Context) error {
	id := c.Param("id")

	var notification model.Notification
	if err := c.Bind(&notification); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid notification ID"))
	}
	notification.ID = objectID

	updatedNotification, err := handler.notificationService.UpdateNotification(&notification)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Notification updated successfully",
		"notification": updatedNotification,
	}, nil))
}

// DeleteNotification handles deleting a notification by ID
func (handler *NotificationHandler) DeleteNotification(c echo.Context) error {
	id := c.Param("id")

	if err := handler.notificationService.DeleteNotification(id); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Notification deleted successfully",
	}, nil))
}

// GetNotificationsByUserID handles retrieving notifications by user ID
func (handler *NotificationHandler) GetNotificationsByUserID(c echo.Context) error {
	userID := c.Param("userID")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Missing user ID"))
	}

	notifications, err := handler.notificationService.GetNotificationsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"notifications": notifications,
	}, nil))
}

// GetNotifications handles retrieving all notifications
func (handler *NotificationHandler) GetNotifications(c echo.Context) error {
	notifications, err := handler.notificationService.GetNotifications()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"notifications": notifications,
	}, nil))
}


