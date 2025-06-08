package di

import (
	"os"

	"github.com/AliAmjid/newsletter-go/internal/db"
	"github.com/AliAmjid/newsletter-go/internal/repository/postgres"
	authusecase "github.com/AliAmjid/newsletter-go/internal/usecase/auth"
	postusecase "github.com/AliAmjid/newsletter-go/internal/usecase/post"
	userusecase "github.com/AliAmjid/newsletter-go/internal/usecase/user"
)

// Container holds dependencies for the application.
type Container struct {
	PostService *postusecase.Service
	AuthService *authusecase.Service
	UserService *userusecase.Service
}

// NewContainer initializes the application dependencies.
func NewContainer() *Container {
	conn := os.Getenv("POSTGRES_CONNECTION_STRING")
	db.Init(conn)

	repo := postgres.NewPostRepository(db.DB)
	service := postusecase.NewService(repo)

	userRepo := postgres.NewUserRepository(db.DB)
	authApiKey := os.Getenv("PERMIT_API_KEY")
	secret := os.Getenv("JWT_SECRET")
	authService := authusecase.NewService(userRepo, authApiKey, []byte(secret))
	userService := userusecase.NewService(userRepo, authApiKey, []byte(secret))

	return &Container{
		PostService: service,
		AuthService: authService,
		UserService: userService,
	}
}
