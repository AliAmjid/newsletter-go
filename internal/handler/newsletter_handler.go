package handler

import (
	"encoding/json"
	"net/http"
	"newsletter-go/domain"
	"newsletter-go/internal/service"

	"github.com/go-chi/chi/v5"
)

type NewsletterHandler struct {
	service *service.NewsletterService
}

func NewNewsletterHandler(s *service.NewsletterService) *NewsletterHandler {
	return &NewsletterHandler{service: s}
}

// GET /newsletters?ownerId=xxx
func (h *NewsletterHandler) List(w http.ResponseWriter, r *http.Request) {
	ownerID := r.URL.Query().Get("ownerId")
	if ownerID == "" {
		http.Error(w, "missing ownerId query param", http.StatusBadRequest)
		return
	}

	newsletters, err := h.service.ListByOwner(r.Context(), ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, newsletters)
}

// POST /newsletters
func (h *NewsletterHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.Newsletter
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, req)
}

// GET /newsletters/{id}
func (h *NewsletterHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	newsletter, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if newsletter == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, newsletter)
}

// PATCH /newsletters/{id}
func (h *NewsletterHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req domain.Newsletter
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.ID = id

	if err := h.service.Update(r.Context(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, req)
}

// DELETE /newsletters/{id}
func (h *NewsletterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}
