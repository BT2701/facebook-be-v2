package routes

import (
	"snake_api/controllers"

	"github.com/labstack/echo/v4"
)

func HealthRoutes(e *echo.Echo) {
	e.GET("/health", controllers.HealthCheck)
}
