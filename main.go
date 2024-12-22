package main

import (
	"log"
	"os"
	"snake_api/config"
	"snake_api/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	config.ConnectDB()

	// Get the user collection from the connected DB
	userCollection := config.GetCollection("users")

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
