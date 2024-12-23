package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/home", func(c echo.Context) error {
		return c.String(http.StatusOK, "Message sent successfully!")
	})

	e.Logger.Fatal(e.Start(":8081"))
}
