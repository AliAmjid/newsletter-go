package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	delivery "newsletter-go/internal/delivery/http"
	"newsletter-go/internal/di"
)

func main() {
	// Load env values
	if err := godotenv.Load(); err != nil {
		log.Fatal("An error has occured during loading from .env file")
	}

	c := di.NewContainer()

	r := delivery.NewRouter()
	//r.Use(middleware.Logger)

	delivery.NewPostHandler(r, c.PostService)
	delivery.NewHelloHandler(r, c.UserService)
	delivery.NewAuthHandler(r, c.AuthService, c.UserService)
	delivery.NewNewsletterHandler(r, c.NewsletterService, c.UserService)

	delivery.NewSubscriberHandler(r, c.SubscriberService)

	fmt.Println("Server starting on port 3000")
	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("failed to listen to server", err)
	}
}
