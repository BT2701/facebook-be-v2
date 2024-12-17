package routes

import (
	"snake_api/controllers"

	"github.com/gin-gonic/gin"
)

func HealthRoutes(router *gin.Engine) {
	router.GET("/health", controllers.HealthCheck)
}
