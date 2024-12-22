package utils

import (
	"github.com/rs/cors"
	"github.com/labstack/echo/v4"
)

func CorsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cors.New(cors.Options{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
				AllowedHeaders: []string{"Origin", "Content-Type"},
			}).HandlerFunc(c.Response().Writer, c.Request())
			return next(c)
		}
	}
}
