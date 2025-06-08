package user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/enforcement"
	permitpkg "github.com/permitio/permit-golang/pkg/permit"

	"github.com/AliAmjid/newsletter-go/domain"
)

type Service struct {
	repo   domain.UserRepository
	permit *permitpkg.Client
	jwtKey []byte
}

func NewService(r domain.UserRepository, apiKey string, jwtKey []byte) *Service {
	cfg := config.NewConfigBuilder(apiKey).WithPdpUrl("https://cloudpdp.api.permit.io").Build()
	return &Service{repo: r, permit: permitpkg.NewPermit(cfg), jwtKey: jwtKey}
}

func (s *Service) tokenFromRequest(r *http.Request) (string, error) {
	h := r.Header.Get("Authentication")
	if h == "" {
		return "", errors.New("missing auth header")
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid auth header")
	}
	return parts[1], nil
}

func (s *Service) parseToken(tokenStr string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil || !t.Valid {
		return "", errors.New("invalid token")
	}
	return claims.Subject, nil
}

func (s *Service) IsLoggedIn(r *http.Request) (*domain.User, error) {
	token, err := s.tokenFromRequest(r)
	if err != nil {
		return nil, err
	}
	userID, err := s.parseToken(token)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(r.Context(), userID)
}

func (s *Service) IsAllowedTo(r *http.Request, action, resource string) (bool, error) {
	token, err := s.tokenFromRequest(r)
	if err != nil {
		return false, err
	}
	userID, err := s.parseToken(token)
	if err != nil {
		return false, err
	}
	u := enforcement.UserBuilder(userID).Build()
	res := enforcement.ResourceBuilder(resource).Build()
	return s.permit.Check(u, enforcement.Action(action), res)
}
