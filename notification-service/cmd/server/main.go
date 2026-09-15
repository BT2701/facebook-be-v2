package main

import (
	"log"
	"notification-service/cmd/route"

	"github.com/BT2701/facebook-be-v2/shared/config"
	"github.com/BT2701/facebook-be-v2/shared/httpx"
)

func main() {
	config.Load()

	router := route.SetupRouter()
	port := config.Get("PORT", "8081")
	log.Printf("notification-service listening on :%s", port)
	if err := httpx.Run(router, port); err != nil {
		log.Fatal(err)
	}
}
