package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BT2701/snake/config"
	"github.com/BT2701/snake/routes"
	"github.com/BT2701/snake/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Load biến môi trường
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Gán JWT Secret
	utils.JWTSecret = []byte(os.Getenv("JWT_SECRET"))

	// Kết nối MongoDB
	config.ConnectDB()

	// Khởi tạo routes
	routes.RegisterRoutes()

	// Khởi chạy server
	port := os.Getenv("PORT")
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
