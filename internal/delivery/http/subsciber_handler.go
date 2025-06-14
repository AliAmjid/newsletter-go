package http

import (
	"encoding/json"
	"net/http"
	"newsletter-go/domain"
	subscriberusecase "newsletter-go/internal/usecase/subscriber"

	"github.com/go-chi/chi/v5"
)

type SubscriberHandler struct {
	service *subscriberusecase.Service
}

func NewSubscriberHandler(r chi.Router, s *subscriberusecase.Service) {
	h := &SubscriberHandler{service: s}

	r.Route("/subscribers", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetByID)
	})
}

func (h *SubscriberHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.User
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

func (h *SubscriberHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, user)
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}
