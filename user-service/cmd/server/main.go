package main

import (
	"log"
	"user-service/cmd/routes"
	"user-service/pkg/database"

	"github.com/BT2701/facebook-be-v2/shared/config"
	"github.com/BT2701/facebook-be-v2/shared/httpx"
)

func main() {
	config.Load()

	database.ConnectDB()
	userCollection := database.GetCollection("users")
	router := routes.SetupRouter(userCollection)

	port := config.Get("PORT", "8080")
	log.Printf("user-service listening on :%s", port)
	if err := httpx.Run(router, port); err != nil {
		log.Fatal(err)
	}
}
