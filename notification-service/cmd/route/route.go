package route

import (
	"notification-service/internal/adapters/inbound"
	"notification-service/internal/adapters/outbound"
	"notification-service/internal/app/services"
	"notification-service/pkg/database"
	"os"

	"github.com/BT2701/facebook-be-v2/shared/httpx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter() *echo.Echo {
	// Initialize MongoDB connection
	database.InitMongoDB()
	databaseName := os.Getenv("DB_NAME")
	notificationCollection := database.GetCollection(databaseName, "notifications")

	// Create repositories and services
	notificationRepo := outbound.NewNotificationRepository(notificationCollection)
	notificationService := services.NewNotificationService(notificationRepo)


	// Create handlers
	notificationHandler := inbound.NewNotificationHandler(notificationService)

	// Set up Echo
	e := echo.New()
	e.Use(middleware.Logger())
	httpx.Apply(e)
	e.POST("/notifications", notificationHandler.CreateNotification)
	e.GET("/notifications/:id", notificationHandler.GetNotification)
	e.PUT("/notifications/:id", notificationHandler.UpdateNotification)
	e.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
	e.GET("/notifications", notificationHandler.GetNotifications)
	e.GET("/notifications/:userID/notifications", notificationHandler.GetNotificationsByUserID)

	return e
}