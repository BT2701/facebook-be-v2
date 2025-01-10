package main

import (
	"chat-service/cmd/route"
	"chat-service/pkg/database"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Initialize MongoDB connection
	database.InitMongoDB()

	// Set up router
	router := route.SetupRouter()

	// Start the server
	log.Println("Starting server...")
	if err := router.Start(os.Getenv(":" + os.Getenv("PORT"))); err != nil {
		log.Fatal("Server error:", err)
	}
}
