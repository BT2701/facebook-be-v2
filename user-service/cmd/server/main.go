package main

import (
	"log"
	"os"
	"user-service/pkg/database"
	"user-service/cmd/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	database.ConnectDB()

	// Get the user collection from the connected DB
	userCollection := database.GetCollection("users")

	// Set up router and pass the collection to SetupRouter
	router := routes.SetupRouter(userCollection)

	// Start server
	port := os.Getenv("PORT")
	log.Printf("Server running on port %s", port)
	err = router.Start(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
