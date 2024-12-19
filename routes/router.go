package routes

import (
	"snake_api/controllers"
	"snake_api/repositories"
	"snake_api/services"
	"snake_api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userCollection *mongo.Collection) *gin.Engine {
	// Khởi tạo repository, service, và controller
	userRepo := repositories.NewUserRepository(userCollection)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})
	r.Use(utils.CorsMiddleware())

	// Routes
	api := r.Group("/api")
	{
		api.POST("/login", userController.Login)
		api.POST("/register", userController.SignUp)
		api.POST("/forgot", userController.ForgotPassword)
		api.POST("/reset", userController.ResetPassword)
	}

	return r
}
