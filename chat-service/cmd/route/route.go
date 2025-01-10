package route

import (
	"chat-service/internal/adapters/in"
	"chat-service/internal/adapters/out"
	"chat-service/internal/app/service"
	"chat-service/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter() *echo.Echo {
	messageRepo := out.NewMongoMessageRepository()
	chatService := service.NewChatService(*messageRepo)

	// Create handlers
	chatHandler := in.ChatHandler{Service: chatService}
	socketHandler := in.NewSocketHandler(chatService)

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

	// HTTP Routes
	api := e.Group("/api")
	api.POST("/messages", chatHandler.SendMessage)
	api.GET("/messages", chatHandler.GetMessages)
	api.GET("/messages/all", chatHandler.GetAllMessages)
	api.DELETE("/messages/all", chatHandler.DeleteAllMessages)

	// WebSocket Route
	api.GET("/ws", socketHandler.HandleConnection)

	return e
}
