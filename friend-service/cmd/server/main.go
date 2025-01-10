package main

import (
	"log"
	"os"
	"friend-service/cmd/route"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	route := route.SetupRouter()

	// Start the server
	if err := route.Start(":" + os.Getenv("PORT")); err != nil {
		log.Fatal("Server error:", err)
	}
}

