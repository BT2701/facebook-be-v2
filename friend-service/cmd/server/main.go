package main

import (
	"log"
	"os"
	"friend-service/internal/adapters/inbound"
	"friend-service/internal/adapters/outbound"
	"friend-service/internal/app/service"
	"friend-service/pkg/database"
	"friend-service/pkg/utils"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}
	// Initialize MongoDB connection
	database.InitMongoDB()
	databaseName := os.Getenv("DB_NAME")
	friendCollection := database.GetCollection(databaseName, "friends")

	// Create repositories and services
	friendRepo := outbound.NewFriendRepository(friendCollection)
	friendService := service.NewFriendService(friendRepo)


	// Create handlers
	friendHandler := inbound.NewFriendHandler(friendService)

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
	e.POST("/friends", friendHandler.CreateFriend)
	e.GET("/friends/:id", friendHandler.GetFriend)
	e.PUT("/friends/:id", friendHandler.UpdateFriend)
	e.DELETE("/friends/:id", friendHandler.DeleteFriend)
	e.GET("/friends", friendHandler.GetFriends)
	e.GET("/friends/:userID/friends", friendHandler.GetFriendsByUserID)
	e.GET("/friends/:userID1/:userID2", friendHandler.IsFriend)

	// Start the server
	if err := e.Start(":" + os.Getenv("PORT")); err != nil {
		log.Fatal("Server error:", err)
	}
}

