package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// HomeHandler xử lý route "/"
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to your Go server!"))
}

func main() {
	// Load config từ file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Lấy port từ .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Port mặc định nếu không có config
	}

	// Khởi tạo router
	router := mux.NewRouter()

	// Định nghĩa route
	router.HandleFunc("/", HomeHandler).Methods("GET")

	// CORS Configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Cho phép tất cả origin, bạn có thể giới hạn cụ thể
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Start server
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler.Handler(router)))
}
