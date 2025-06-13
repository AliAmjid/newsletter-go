package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	firebase "firebase.google.com/go/v4"
	fbauth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"github.com/google/uuid"
	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/models"
    permitpkg "github.com/permitio/permit-golang/pkg/permit"
    "newsletter-go/internal/mailer"

	"newsletter-go/domain"
)

type Service struct {
	repo         domain.UserRepository
	resetRepo    domain.PasswordResetRepository
	permit       *permitpkg.Client
	firebaseKey  string
	authClient   *fbauth.Client
    mailer       *mailer.Service
}

type signUpResponse struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	LocalID      string `json:"localId"`
}

func NewService(r domain.UserRepository, rr domain.PasswordResetRepository, permitKey, creds string, firebaseKey string, m *mailer.Service) *Service {
	cfg := config.NewConfigBuilder(permitKey).WithPdpUrl("https://cloudpdp.api.permit.io").Build()
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(creds))
	if err != nil {
		panic(err)
	}
	ac, err := app.Auth(context.Background())
	if err != nil {
		panic(err)
	}
    return &Service{repo: r, resetRepo: rr, permit: permitpkg.NewPermit(cfg), firebaseKey: firebaseKey, authClient: ac, mailer: m}
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
	u := &domain.User{Email: email, FirebaseUID: res.LocalID}
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
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil || u == nil {
		return fmt.Errorf("user not found")
	}

	token := uuid.New().String()
	expires := time.Now().Add(time.Hour).Unix()
	if err := s.resetRepo.Create(ctx, &domain.PasswordResetToken{Token: token, UserID: u.ID, ExpiresAt: expires}); err != nil {
		return err
	}

    return s.mailer.SendForgotPasswordEmail(u.Email, token)
}

func (s *Service) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	rt, err := s.resetRepo.Get(ctx, token)
	if err != nil || rt == nil {
		return fmt.Errorf("invalid token")
	}
	if time.Now().Unix() > rt.ExpiresAt {
		return fmt.Errorf("token expired")
	}

	if _, err := s.authClient.UpdateUser(ctx, rt.UserID, (&fbauth.UserToUpdate{}).Password(newPassword)); err != nil {
		return err
	}
	if err := s.resetRepo.Delete(ctx, token); err != nil {
		return err
	}
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

type firebaseError struct {
	Error struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
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

	fmt.Printf("Response body: %s,%s\n", string(req.RemoteAddr), string(bodyBytes))

	if resp.StatusCode >= 300 {
		var fbErr firebaseError
		if err := json.Unmarshal(bodyBytes, &fbErr); err != nil {
			return fmt.Errorf("firebase error: %s", resp.Status)
		}
		return fmt.Errorf("firebase error: %s", fbErr.Error.Message)
	}

	return json.Unmarshal(bodyBytes, out)
}
