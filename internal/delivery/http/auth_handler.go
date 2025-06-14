package http

import (
	"net/http"
	userusecase "newsletter-go/internal/usecase/user"

	"github.com/go-chi/chi/v5"
	authusecase "newsletter-go/internal/usecase/auth"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service  *authusecase.Service
	users    *userusecase.Service
	validate *validator.Validate
}

func NewAuthHandler(r chi.Router, s *authusecase.Service, u *userusecase.Service) {
	h := &AuthHandler{
		service:  s,
		users:    u,
		validate: validator.New(),
	}
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
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type passwordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type passwordResetConfirm struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func (h *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if !bindAndValidate(w, r, &req, h.validate) {
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
	if !bindAndValidate(w, r, &req, h.validate) {
		return
	}
	access, refresh, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == authusecase.ErrInvalidCredentials {
			respondWithError(w, http.StatusUnauthorized, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "login failed")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"accessToken": access, "refreshToken": refresh})
}

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if !bindAndValidate(w, r, &req, h.validate) {
		return
	}
	access, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if err == authusecase.ErrInvalidToken {
			respondWithError(w, http.StatusUnauthorized, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "failed to refresh token")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"accessToken": access})
}

func (h *AuthHandler) requestReset(w http.ResponseWriter, r *http.Request) {
	var req passwordResetRequest
	if !bindAndValidate(w, r, &req, h.validate) {
		return
	}
	if err := h.service.RequestPasswordReset(r.Context(), req.Email); err != nil {
		if err == authusecase.ErrUserNotFound {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "failed to request reset")
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) confirmReset(w http.ResponseWriter, r *http.Request) {
	var req passwordResetConfirm
	if !bindAndValidate(w, r, &req, h.validate) {
		return
	}
	if err := h.service.ConfirmPasswordReset(r.Context(), req.Token, req.NewPassword); err != nil {
		switch err {
		case authusecase.ErrInvalidToken:
			respondWithError(w, http.StatusBadRequest, err.Error())
		case authusecase.ErrTokenExpired:
			respondWithError(w, http.StatusBadRequest, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to reset password")
		}
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
	respondWithJSON(w, http.StatusOK, map[string]string{"id": u.ID, "email": u.Email})
}
