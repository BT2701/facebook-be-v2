package routes

import (
	"net/http"

	"github.com/BT2701/snake/controllers"
)

func RegisterRoutes() {
	http.HandleFunc("/register", controllers.RegisterHandler)
	http.HandleFunc("/login", controllers.LoginHandler)
}
