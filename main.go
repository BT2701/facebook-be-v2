package main

import (
	"log"
	"os"
	"snake_api/config"
	"snake_api/controllers"
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
	controllers.InitUserController(config.DB.Collection("users"))


	// Set up router
	router := routes.SetupRouter()

	// Start server
	port := os.Getenv("PORT")
	log.Printf("Server running on port %s", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}