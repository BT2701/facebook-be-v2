package route

import (
	"os"
	"notification-service/internal/adapters/inbound"
	"notification-service/internal/adapters/outbound"
	"notification-service/internal/app/services"
	"notification-service/pkg/database"
	"notification-service/pkg/utils"

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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			return next(c)
		}
	})
	e.Use(utils.CorsMiddleware())
	e.POST("/notifications", notificationHandler.CreateNotification)
	e.GET("/notifications/:id", notificationHandler.GetNotification)
	e.PUT("/notifications/:id", notificationHandler.UpdateNotification)
	e.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
	e.GET("/notifications", notificationHandler.GetNotifications)
	e.GET("/notifications/:userID/notifications", notificationHandler.GetNotificationsByUserID)

	return e
}