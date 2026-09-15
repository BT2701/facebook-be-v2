package main

import (
	"chat-service/cmd/route"
	"chat-service/pkg/database"
	"log"

	"github.com/BT2701/facebook-be-v2/shared/config"
	"github.com/BT2701/facebook-be-v2/shared/httpx"
)

func main() {
	config.Load()
	database.InitMongoDB()

	router := route.SetupRouter()
	port := config.Get("PORT", "8082")
	log.Printf("chat-service listening on :%s", port)
	if err := httpx.Run(router, port); err != nil {
		log.Fatal(err)
	}
}
