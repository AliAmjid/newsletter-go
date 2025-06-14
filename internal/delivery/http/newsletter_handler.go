package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"newsletter-go/domain"
	newsletterusecase "newsletter-go/internal/usecase/newsletter"
	userusecase "newsletter-go/internal/usecase/user"
)

// NewsletterHandler handles newsletter HTTP endpoints.
type NewsletterHandler struct {
	service  *newsletterusecase.Service
	users    *userusecase.Service
	validate *validator.Validate
}

func NewNewsletterHandler(r chi.Router, s *newsletterusecase.Service, u *userusecase.Service) {
	h := &NewsletterHandler{service: s, users: u, validate: validator.New()}

	r.Route("/newsletters", func(r chi.Router) {
		r.Get("/", h.listNewsletters)
		r.Post("/", h.createNewsletter)
		r.Get("/{newsletterId}", h.getNewsletter)
		r.Patch("/{newsletterId}", h.updateNewsletter)
		r.Delete("/{newsletterId}", h.deleteNewsletter)
	})
}

type newsletterCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type newsletterUpdateRequest struct {
	Name        string `json:"name" validate:"omitempty"`
	Description string `json:"description"`
}

func (h *NewsletterHandler) listNewsletters(w http.ResponseWriter, r *http.Request) {
	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if ok, err := h.users.IsAllowedTo(r, "read", "newsletter"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	list, err := h.service.ListByOwner(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list newsletters")
		return
	}
	respondWithJSON(w, http.StatusOK, list)
}

func (h *NewsletterHandler) createNewsletter(w http.ResponseWriter, r *http.Request) {
	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if ok, err := h.users.IsAllowedTo(r, "create", "newsletter"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req newsletterCreateRequest
	if !bindAndValidate(w, r, &req, h.validate) {
		return
	}

	n := &domain.Newsletter{Name: req.Name, Description: req.Description, OwnerID: user.ID}
	if err := h.service.Create(r.Context(), n); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create newsletter")
		return
	}
	respondWithJSON(w, http.StatusCreated, n)
}

func (h *NewsletterHandler) getNewsletter(w http.ResponseWriter, r *http.Request) {
	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if ok, err := h.users.IsAllowedTo(r, "read", "newsletter"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	id := chi.URLParam(r, "newsletterId")
	n, err := h.service.GetByID(r.Context(), id)
	if err != nil || n == nil {
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}
	if n.OwnerID != user.ID {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}
	respondWithJSON(w, http.StatusOK, n)
}

func (h *NewsletterHandler) updateNewsletter(w http.ResponseWriter, r *http.Request) {
	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if ok, err := h.users.IsAllowedTo(r, "update", "newsletter"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	id := chi.URLParam(r, "newsletterId")
	n, err := h.service.GetByID(r.Context(), id)
	if err != nil || n == nil {
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}
	if n.OwnerID != user.ID {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req newsletterUpdateRequest
	if !bindAndValidate(w, r, &req, h.validate) {
		return
	}

	if req.Name != "" {
		n.Name = req.Name
	}
	if req.Description != "" {
		n.Description = req.Description
	}
	if err := h.service.Update(r.Context(), n); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update newsletter")
		return
	}
	respondWithJSON(w, http.StatusOK, n)
}

func (h *NewsletterHandler) deleteNewsletter(w http.ResponseWriter, r *http.Request) {
	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if ok, err := h.users.IsAllowedTo(r, "delete", "newsletter"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	id := chi.URLParam(r, "newsletterId")
	n, err := h.service.GetByID(r.Context(), id)
	if err != nil || n == nil {
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}
	if n.OwnerID != user.ID {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete newsletter")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
