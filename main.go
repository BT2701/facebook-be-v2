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

	// Set up router
	router := routes.SetupRouter()
	routes.HealthRoutes(router)
	routes.AuthRoutes(router)


	// Start server
	port := os.Getenv("PORT")
	log.Printf("Server running on port %s", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
