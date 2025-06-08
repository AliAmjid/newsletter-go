package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/models"
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

func (s *Service) SignUp(ctx context.Context, email, password string) (string, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	u := &domain.User{Email: email, PasswordHash: string(hash)}
	if err := s.repo.Create(ctx, u); err != nil {
		return "", "", err
	}
	newUser := models.NewUserCreate(u.ID)
	newUser.Email = &u.Email
	if _, err := s.permit.SyncUser(ctx, *newUser); err != nil {
		return "", "", err
	}
	return s.issueAccessToken(u.ID), s.issueRefreshToken(u.ID), nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, string, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil || u == nil {
		return "", "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", "", errors.New("invalid credentials")
	}
	return s.issueAccessToken(u.ID), s.issueRefreshToken(u.ID), nil
}

func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	// In real implementation we would send email with token
	// For now we just return nil
	return nil
}

func (s *Service) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	// This example does not store reset tokens, so just a stub
	return nil
}

func (s *Service) issueAccessToken(userID string) string {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, _ := token.SignedString(s.jwtKey)
	return str
}

func (s *Service) issueRefreshToken(userID string) string {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, _ := token.SignedString(s.jwtKey)
	return str
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

func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	userID, err := s.parseToken(refreshToken)
	if err != nil {
		return "", err
	}
	return s.issueAccessToken(userID), nil
}
