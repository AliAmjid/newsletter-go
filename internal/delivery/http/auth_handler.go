package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	authusecase "github.com/AliAmjid/newsletter-go/internal/usecase/auth"
)

type AuthHandler struct {
	service *authusecase.Service
}

func NewAuthHandler(r chi.Router, s *authusecase.Service) {
	h := &AuthHandler{service: s}
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", h.signUp)
		r.Post("/login", h.login)
		r.Post("/password-reset/request", h.requestReset)
		r.Post("/password-reset/confirm", h.confirmReset)
	})
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type passwordResetRequest struct {
	Email string `json:"email"`
}

type passwordResetConfirm struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

func (h *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	token, err := h.service.SignUp(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) requestReset(w http.ResponseWriter, r *http.Request) {
	var req passwordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if err := h.service.RequestPasswordReset(r.Context(), req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) confirmReset(w http.ResponseWriter, r *http.Request) {
	var req passwordResetConfirm
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if err := h.service.ConfirmPasswordReset(r.Context(), req.Token, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
