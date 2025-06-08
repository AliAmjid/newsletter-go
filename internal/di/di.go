package di

import (
	"os"

	"newsletter-go/internal/db"
	"newsletter-go/internal/repository/postgres"
	authusecase "newsletter-go/internal/usecase/auth"
	postusecase "newsletter-go/internal/usecase/post"
	userusecase "newsletter-go/internal/usecase/user"
)

// Container holds dependencies for the application.
type Container struct {
	PostService *postusecase.Service
	AuthService *authusecase.Service
	UserService *userusecase.Service
}

func NewContainer() *Container {
	conn := os.Getenv("POSTGRES_CONNECTION_STRING")
	db.Init(conn)

	repo := postgres.NewPostRepository(db.DB)
	service := postusecase.NewService(repo)

	userRepo := postgres.NewUserRepository(db.DB)
	resetRepo := postgres.NewPasswordResetRepository(db.DB)
	authApiKey := os.Getenv("PERMIT_API_KEY")
	fbCreds := os.Getenv("FIREBASE_CREDENTIALS")
	fbKey := os.Getenv("FIREBASE_API_KEY")
	sgKey := os.Getenv("SENDGRID_API_KEY")
	sgFrom := os.Getenv("SENDGRID_FROM_EMAIL")
	authService := authusecase.NewService(userRepo, resetRepo, authApiKey, fbCreds, fbKey, sgKey, sgFrom)
	userService := userusecase.NewService(userRepo, authApiKey, fbCreds)

	return &Container{
		PostService: service,
		AuthService: authService,
		UserService: userService,
	}
}
