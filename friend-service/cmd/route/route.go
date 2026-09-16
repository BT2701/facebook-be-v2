package route

import (
	"friend-service/internal/adapters/inbound"
	"friend-service/internal/adapters/outbound"
	"friend-service/internal/app/service"
	"friend-service/pkg/database"
	"os"

	"github.com/BT2701/facebook-be-v2/shared/events"
	"github.com/BT2701/facebook-be-v2/shared/httpx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter() *echo.Echo {
	database.InitMongoDB()
	databaseName := os.Getenv("DB_NAME")
	friendCollection := database.GetCollection(databaseName, "friends")
	requestCollection := database.GetCollection(databaseName, "requests")

	friendRepo := outbound.NewFriendRepository(friendCollection)
	friendService := service.NewFriendService(friendRepo)
	requestRepo := outbound.NewRequestRepository(requestCollection)
	requestService := service.NewRequestService(requestRepo)

	bus := events.New(os.Getenv("REDIS_URI"))
	friendHandler := inbound.NewFriendHandlerWithRequests(friendService, requestService, bus)
	requestHandler := inbound.NewRequestHandler(requestService, bus)

	e := echo.New()
	e.Use(middleware.Logger())
	httpx.Apply(e)

	e.POST("/friends", friendHandler.CreateFriend)
	e.POST("/", friendHandler.CreateFriend)
	e.GET("/friends/:userID/friends", friendHandler.GetFriendsByUserID)
	e.GET("/friends/:userID1/:userID2", friendHandler.IsFriend)
	e.GET("/friends", friendHandler.GetFriends)
	e.DELETE("/friends/:id", friendHandler.DeleteFriend)
	e.DELETE("/remove/:userId/:friendId", friendHandler.RemoveFriend)
	e.POST("/create-and-delete-request", friendHandler.CreateAndDeleteRequest)
	e.GET("/nonfriends/:userID", friendHandler.GetNonFriends)

	e.POST("/Request", requestHandler.CreateRequest)
	e.GET("/Request", requestHandler.GetRequestsByQuery)
	e.GET("/Request/requests", requestHandler.GetRequestsByQuery)
	e.GET("/Request/:sender/:receiver", requestHandler.GetRequest)
	e.DELETE("/Request/delete", requestHandler.DeleteRequest)
	e.DELETE("/Request/:id", requestHandler.DeleteRequestByID)

	e.POST("/requests", requestHandler.CreateRequest)
	e.GET("/requests", requestHandler.GetRequestsByQuery)
	e.GET("/requests/:sender/:receiver", requestHandler.GetRequest)
	e.DELETE("/requests/:sender/:receiver", requestHandler.DeleteRequest)

	return e
}
