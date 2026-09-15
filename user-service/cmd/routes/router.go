package routes

import (
	"os"
	"user-service/internal/adapters/inbound"
	"user-service/internal/adapters/outbound"
	"user-service/internal/app/services"

	"github.com/BT2701/facebook-be-v2/shared/httpx"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(userCollection *mongo.Collection) *echo.Echo {
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URI"), // Địa chỉ Redis
	})
	// Khởi tạo repository, service, và controller
	userRepo := outbound.NewUserRepository(userCollection)
	userService := services.NewUserService(userRepo, redisClient)
	userController := inbound.NewUserController(userService)

	e := echo.New()
	httpx.Apply(e)

	api := e.Group("/api")
	{
		api.POST("/login", userController.Login)
		api.POST("/register", userController.SignUp)
		api.POST("/forgot", userController.ForgotPassword)
		api.POST("/reset", userController.ResetPassword)
		api.GET("/users", userController.GetAllUsers)
		api.PUT("/logout", userController.Logout)
		api.PUT("/edit", userController.EditUser)
		api.GET("/user/:id", userController.GetByID)
		api.GET("/user", userController.FindUserByEmail)
		api.PUT("/avatar", userController.UpdateAvatar)
	}

	return e
}
