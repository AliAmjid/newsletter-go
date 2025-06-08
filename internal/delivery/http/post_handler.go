package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AliAmjid/newsletter-go/domain"
	postusecase "github.com/AliAmjid/newsletter-go/internal/usecase/post"
)

type PostHandler struct {
	service *postusecase.Service
}

// NewPostHandler registers post endpoints to the router.
func NewPostHandler(r chi.Router, s *postusecase.Service) {
	h := &PostHandler{service: s}

	r.Route("/posts", func(r chi.Router) {
		r.Post("/", h.createPost)
		r.Get("/{id}", h.getPost)
	})
}

func (h *PostHandler) createPost(w http.ResponseWriter, r *http.Request) {
	var p domain.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.Save(r.Context(), &p); err != nil {
		http.Error(w, "Failed to save post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Post recieved"))
}

func (h *PostHandler) getPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("TODO post with ID " + id))
}
