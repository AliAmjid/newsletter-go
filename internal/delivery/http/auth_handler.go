package http

import (
	"encoding/json"
	"net/http"
	userusecase "newsletter-go/internal/usecase/user"

	"github.com/go-chi/chi/v5"

	authusecase "newsletter-go/internal/usecase/auth"
)

type AuthHandler struct {
	service *authusecase.Service
	users   *userusecase.Service
}

func NewAuthHandler(r chi.Router, s *authusecase.Service, u *userusecase.Service) {
	h := &AuthHandler{service: s, users: u}
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", h.signUp)
		r.Post("/login", h.login)
		r.Post("/refresh", h.refresh)
		r.Post("/password-reset/request", h.requestReset)
		r.Post("/password-reset/confirm", h.confirmReset)
		r.Get("/whoami", h.whoAmI)
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
		respondWithError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	access, refresh, err := h.service.SignUp(r.Context(), req.Email, req.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"accessToken": access, "refreshToken": refresh})
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	access, refresh, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"accessToken": access, "refreshToken": refresh})
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	access, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"accessToken": access})
}

func (h *AuthHandler) requestReset(w http.ResponseWriter, r *http.Request) {
	var req passwordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := h.service.RequestPasswordReset(r.Context(), req.Email); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) confirmReset(w http.ResponseWriter, r *http.Request) {
	var req passwordResetConfirm
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := h.service.ConfirmPasswordReset(r.Context(), req.Token, req.NewPassword); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) whoAmI(w http.ResponseWriter, r *http.Request) {
	u, err := h.users.IsLoggedIn(r)
	if err != nil || u == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}
