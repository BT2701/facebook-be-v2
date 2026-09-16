package route

import (
	"context"
	"log"
	"notification-service/internal/adapters/inbound"
	"notification-service/internal/adapters/outbound"
	"notification-service/internal/app/services"
	"notification-service/pkg/database"
	"os"

	"github.com/BT2701/facebook-be-v2/shared/events"
	"github.com/BT2701/facebook-be-v2/shared/httpx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter() *echo.Echo {
	database.InitMongoDB()
	databaseName := os.Getenv("DB_NAME")
	notificationCollection := database.GetCollection(databaseName, "notifications")

	notificationRepo := outbound.NewNotificationRepository(notificationCollection)
	notificationService := services.NewNotificationService(notificationRepo)
	notificationHandler := inbound.NewNotificationHandler(notificationService)

	bus := events.New(os.Getenv("REDIS_URI"))
	go bus.Subscribe(context.Background(), func(event events.NotificationEvent) {
		if err := notificationService.ApplyEvent(event); err != nil {
			log.Printf("notification event: %v", err)
		}
	}, func(event events.NotificationEvent) {
		if err := notificationService.ApplyDeleteEvent(event); err != nil {
			log.Printf("notification delete event: %v", err)
		}
	})

	e := echo.New()
	e.Use(middleware.Logger())
	httpx.Apply(e)

	e.POST("/notifications", notificationHandler.CreateNotification)
	e.POST("/", notificationHandler.CreateNotification)
	e.GET("/notifications/:id", notificationHandler.GetNotification)
	e.PUT("/notifications/:id", notificationHandler.UpdateNotification)
	e.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
	e.GET("/notifications", notificationHandler.GetNotifications)
	e.GET("/notifications/:userID/notifications", notificationHandler.GetNotificationsByUserID)

	e.GET("/receiver/:userID", notificationHandler.GetNotificationsByUserID)
	e.PUT("/:id", notificationHandler.MarkAsRead)
	e.PUT("/markAllAsRead/:userID", notificationHandler.MarkAllAsRead)
	e.DELETE("/delete/:user/:receiver/:post/:action", notificationHandler.DeleteByCombo)

	return e
}
