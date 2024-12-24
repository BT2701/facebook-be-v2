package main

import (
    "log"
    "chat-service/internal/adapters/in"
    "chat-service/internal/adapters/out"
    "chat-service/internal/app/service"
    "chat-service/pkg/database"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
	"os"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Error loading .env file:", err)
    }
    // Initialize MongoDB connection
    database.InitMongoDB()

    // Create repositories and services
    messageRepo := out.NewMongoMessageRepository()
    chatService := service.NewChatService(*messageRepo)

    // Create HTTP handler
    chatHandler := in.ChatHandler{Service: chatService}

    // Set up Echo
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Routes
    api := e.Group("/api")
    api.POST("/messages", chatHandler.SendMessage)
    api.GET("/messages", chatHandler.GetMessages)
    api.GET("/messages/all", chatHandler.GetAllMessages)
    api.DELETE("/messages/all", chatHandler.DeleteAllMessages)

    // Start the server
    if err := e.Start(":" + os.Getenv("PORT")); err != nil {
        log.Fatal("Server error:", err)
    }
}
