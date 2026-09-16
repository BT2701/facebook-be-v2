package route

import (
	"ai-service/internal/adapters/inbound"
	"ai-service/internal/adapters/outbound"
	"ai-service/internal/app/service"
	"ai-service/pkg/database"
	"os"

	"github.com/BT2701/facebook-be-v2/shared/httpx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter() *echo.Echo {
	database.InitMongoDB()
	databaseName := os.Getenv("DB_NAME")
	if databaseName == "" {
		databaseName = "aidb"
	}

	conversationRepo := outbound.NewConversationRepository(database.GetCollection(databaseName, "conversations"))
	messageRepo := outbound.NewMessageRepository(database.GetCollection(databaseName, "messages"))
	aiService := service.NewAIService(conversationRepo, messageRepo, outbound.NewLLMClient())
	handler := inbound.NewAIHandler(aiService)

	e := echo.New()
	e.Use(middleware.Logger())
	httpx.Apply(e)

	e.GET("/chat", handler.GetInbox)
	e.POST("/chat", handler.Chat)
	e.DELETE("/conversations/:id", handler.Clear)

	return e
}
