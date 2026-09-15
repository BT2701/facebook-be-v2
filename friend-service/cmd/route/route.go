package route

import (
	"friend-service/internal/adapters/inbound"
	"friend-service/internal/adapters/outbound"
	"friend-service/internal/app/service"
	"friend-service/pkg/database"
	"os"

	"github.com/BT2701/facebook-be-v2/shared/httpx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter() *echo.Echo {
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
	e.Use(middleware.Logger())
	httpx.Apply(e)
	e.POST("/friends", friendHandler.CreateFriend)
	e.GET("/friends/:id", friendHandler.GetFriend)
	e.PUT("/friends/:id", friendHandler.UpdateFriend)
	e.DELETE("/friends/:id", friendHandler.DeleteFriend)
	e.GET("/friends", friendHandler.GetFriends)
	e.GET("/friends/:userID/friends", friendHandler.GetFriendsByUserID)
	e.GET("/friends/:userID1/:userID2", friendHandler.IsFriend)
	return e
}
