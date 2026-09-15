package main

import (
	"log"
	"media-service/cmd/routes"
	"media-service/pkg/database"

	"github.com/BT2701/facebook-be-v2/shared/config"
	"github.com/BT2701/facebook-be-v2/shared/httpx"
)

func main() {
	config.Load()

	database.ConnectDB()
	imageCollection := database.GetCollection("images")
	router := routes.SetupRouter(imageCollection)

	port := config.Get("PORT", "8083")
	log.Printf("media-service listening on :%s", port)
	if err := httpx.Run(router, port); err != nil {
		log.Fatal(err)
	}
}
