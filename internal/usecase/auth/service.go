package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/models"
	permitpkg "github.com/permitio/permit-golang/pkg/permit"

	"newsletter-go/domain"
)

type Service struct {
	repo        domain.UserRepository
	permit      *permitpkg.Client
	firebaseKey string
}

type signUpResponse struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	LocalID      string `json:"localId"`
}

func NewService(r domain.UserRepository, permitKey, _ string, firebaseKey string) *Service {
	cfg := config.NewConfigBuilder(permitKey).WithPdpUrl("https://cloudpdp.api.permit.io").Build()
	return &Service{repo: r, permit: permitpkg.NewPermit(cfg), firebaseKey: firebaseKey}
}

func (s *Service) firebaseSignUp(ctx context.Context, email, password string) (*signUpResponse, error) {
	body := map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}
	b, _ := json.Marshal(body)
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=%s", s.firebaseKey)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	var resp signUpResponse
	if err := doJSON(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type loginResponse struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	LocalID      string `json:"localId"`
}

func (s *Service) firebaseLogin(ctx context.Context, email, password string) (*loginResponse, error) {
	body := map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}
	b, _ := json.Marshal(body)
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", s.firebaseKey)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	var resp loginResponse
	if err := doJSON(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) SignUp(ctx context.Context, email, password string) (string, string, error) {
	res, err := s.firebaseSignUp(ctx, email, password)
	if err != nil {
		return "", "", err
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := &domain.User{ID: res.LocalID, Email: email, PasswordHash: string(hash)}
	if err := s.repo.Create(ctx, u); err != nil {
		return "", "", err
	}
	newUser := models.NewUserCreate(u.ID)
	newUser.Email = &u.Email
	if _, err := s.permit.SyncUser(ctx, *newUser); err != nil {
		return "", "", err
	}
	return res.IDToken, res.RefreshToken, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, string, error) {
	res, err := s.firebaseLogin(ctx, email, password)
	if err != nil {
		return "", "", err
	}
	return res.IDToken, res.RefreshToken, nil
}

func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	return nil
}

func (s *Service) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	return nil
}

type refreshResponse struct {
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	body := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}
	b, _ := json.Marshal(body)
	url := fmt.Sprintf("https://securetoken.googleapis.com/v1/token?key=%s", s.firebaseKey)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	var resp refreshResponse
	if err := doJSON(req, &resp); err != nil {
		return "", err
	}
	return resp.IDToken, nil
}

func doJSON(req *http.Request, out interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Response body: %s\n", string(bodyBytes))
	if resp.StatusCode >= 300 {
		return fmt.Errorf("firebase error: %s", resp.Status)
	}
	return json.Unmarshal(bodyBytes, out)
}
