package di

import (
	"os"

	"github.com/AliAmjid/newsletter-go/internal/db"
	"github.com/AliAmjid/newsletter-go/internal/repository/postgres"
	postusecase "github.com/AliAmjid/newsletter-go/internal/usecase/post"
)

// Container holds dependencies for the application.
type Container struct {
	PostService *postusecase.Service
}

// NewContainer initializes the application dependencies.
func NewContainer() *Container {
	conn := os.Getenv("POSTGRES_CONNECTION_STRING")
	db.Init(conn)

	repo := postgres.NewPostRepository(db.DB)
	service := postusecase.NewService(repo)

	return &Container{
		PostService: service,
	}
}
