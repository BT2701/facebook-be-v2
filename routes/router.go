package routes

import (
	"snake_api/controllers"
	"snake_api/utils"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	
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
		api.GET("/users", controllers.GetUsers)
		api.POST("/users", controllers.CreateUser)
		api.POST("/login", controllers.Login)
		api.GET("/health", controllers.HealthCheck)
		api.POST("/register", controllers.SignUp)
	}

	return r
}
