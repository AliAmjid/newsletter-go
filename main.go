package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AliAmjid/newsletter-go/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load env values
	err := godotenv.Load()
	if err != nil {
		log.Fatal("An error has occured during loading from .env file")
	}
	// env values can be used with os.Getenv("POSTGRES_HOST") command

	startServer()
}

func startServer() {
	fmt.Println("Server starting on port 3000")

	router := router.NewRouter()

	server := &http.Server{
		Addr:    ":3001",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to listen to server", err)
	}
}
