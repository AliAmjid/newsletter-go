package user

import (
	"errors"
	"net/http"
	"strings"

	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/enforcement"
	permitpkg "github.com/permitio/permit-golang/pkg/permit"
	"google.golang.org/api/option"

	"newsletter-go/domain"
)

type Service struct {
	repo       domain.UserRepository
	permit     *permitpkg.Client
	authClient *auth.Client
}

var (
	ErrMissingAuthHeader = errors.New("missing authorization header")
	ErrEmptyToken        = errors.New("empty token")
	ErrInvalidToken      = errors.New("invalid token")
)

func NewService(r domain.UserRepository, permitKey, creds string) *Service {
	cfg := config.NewConfigBuilder(permitKey).WithPdpUrl("https://cloudpdp.api.permit.io").Build()
	var app *firebase.App
	var err error
	if len(creds) > 0 && creds[0] == '{' {
		app, err = firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON([]byte(creds)))
	} else {
		app, err = firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(creds))
	}
	if err != nil {
		panic(err)
	}
	ac, err := app.Auth(context.Background())
	if err != nil {
		panic(err)
	}
	return &Service{repo: r, permit: permitpkg.NewPermit(cfg), authClient: ac}
}

func (s *Service) tokenFromRequest(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", ErrMissingAuthHeader
	}

	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))
	if token == "" {
		return "", ErrEmptyToken
	}

	return token, nil
}

func (s *Service) parseToken(ctx context.Context, tokenStr string) (string, error) {
	t, err := s.authClient.VerifyIDToken(ctx, tokenStr)
	if err != nil {
		return "", ErrInvalidToken
	}
	return t.UID, nil
}

func (s *Service) IsLoggedIn(r *http.Request) (*domain.User, error) {
	token, err := s.tokenFromRequest(r)
	if err != nil {
		return nil, err
	}
	userID, err := s.parseToken(r.Context(), token)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByFirebaseID(r.Context(), userID)
}

func (s *Service) IsAllowedTo(r *http.Request, action, resource string) (bool, error) {
	user, err := s.IsLoggedIn(r)
	if err != nil {
		return false, err
	}
	u := enforcement.UserBuilder(user.ID).WithEmail(user.Email).Build()
	res := enforcement.ResourceBuilder(resource).Build()
	return s.permit.Check(u, enforcement.Action(action), res)
}
