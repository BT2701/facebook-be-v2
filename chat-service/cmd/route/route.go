package route

import (
	"chat-service/internal/adapters/in"
	"chat-service/internal/adapters/out"
	"chat-service/internal/app/service"

	"github.com/BT2701/facebook-be-v2/shared/httpx"
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
	e.Use(middleware.Logger())
	httpx.Apply(e)

	api := e.Group("/api")
	api.POST("/messages", chatHandler.SendMessage)
	api.GET("/messages", chatHandler.GetMessages)
	api.GET("/messages/all", chatHandler.GetAllMessages)

	// WebSocket Route
	api.GET("/ws", socketHandler.HandleConnection)

	return e
}
