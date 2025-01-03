package main

import (
	"log"
	"os"
	"post-service/internal/adapters/out"
	"post-service/internal/app/database"
	"post-service/internal/app/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"post-service/pkg/utils"
	
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Error loading .env file:", err)
    }
    // Initialize MongoDB connection
    database.InitMongoDB()

    // Create repositories and services
    postRepo := out.NewMongoPostRepository()
	postService := service.NewPostService(*postRepo)




    // Create handlers
    // chatHandler := in.ChatHandler{Service: chatService}
    // socketHandler := in.NewSocketHandler(chatService)

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
    // api := e.Group("/api")
    // api.POST("/messages", chatHandler.SendMessage)
    // api.GET("/messages", chatHandler.GetMessages)
    // api.GET("/messages/all", chatHandler.GetAllMessages)
    // api.DELETE("/messages/all", chatHandler.DeleteAllMessages)

    // WebSocket Route
    // api.GET("/ws", socketHandler.HandleConnection)

    // Start the server
    if err := e.Start(":" + os.Getenv("PORT")); err != nil {
        log.Fatal("Server error:", err)
    }
}
