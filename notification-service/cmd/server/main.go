package main

import (
	"log"
	"os"
	"notification-service/cmd/route"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}
	router := route.SetupRouter()
	// Start the server
	if err := router.Start(":" + os.Getenv("PORT")); err != nil {
		log.Fatal("Server error:", err)
	}
}

