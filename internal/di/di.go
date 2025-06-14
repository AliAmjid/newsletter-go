package di

import (
	"os"

	"newsletter-go/internal/db"
	"newsletter-go/internal/mailer"
	"newsletter-go/internal/repository/postgres"
	authusecase "newsletter-go/internal/usecase/auth"
	newsletterusecase "newsletter-go/internal/usecase/newsletter"
	postusecase "newsletter-go/internal/usecase/post"
	userusecase "newsletter-go/internal/usecase/user"
)

// Container holds dependencies for the application.
type Container struct {
	PostService       *postusecase.Service
	AuthService       *authusecase.Service
	UserService       *userusecase.Service
	NewsletterService *newsletterusecase.Service
}

func NewContainer() *Container {
	conn := os.Getenv("POSTGRES_CONNECTION_STRING")
	db.Init(conn)

	repo := postgres.NewPostRepository(db.DB)
	service := postusecase.NewService(repo)

	newsletterRepo := postgres.NewNewsletterRepository(db.DB)
	newsletterService := newsletterusecase.NewService(newsletterRepo)

	userRepo := postgres.NewUserRepository(db.DB)
	resetRepo := postgres.NewPasswordResetRepository(db.DB)
	authApiKey := os.Getenv("PERMIT_API_KEY")
	fbCreds := os.Getenv("FIREBASE_CREDENTIALS")
	fbKey := os.Getenv("FIREBASE_API_KEY")
	mgDomain := os.Getenv("MAILGUN_DOMAIN")
	mgKey := os.Getenv("MAILGUN_API_KEY")
	mgFrom := os.Getenv("MAILGUN_FROM_EMAIL")
	mailerSvc, err := mailer.NewService(mgDomain, mgKey, mgFrom)
	if err != nil {
		panic(err)
	}
	authService := authusecase.NewService(userRepo, resetRepo, authApiKey, fbCreds, fbKey, mailerSvc)
	userService := userusecase.NewService(userRepo, authApiKey, fbCreds)

	return &Container{
		PostService:       service,
		AuthService:       authService,
		UserService:       userService,
		NewsletterService: newsletterService,
	}
}
