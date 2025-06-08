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

func NewPostHandler(r chi.Router, s *postusecase.Service) {
	h := &PostHandler{service: s}

	r.Route("/newsletters/{newsletterId}/posts", func(r chi.Router) {
		r.Post("/", h.createPost)
		r.Get("/", h.listPosts)
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

func (h *PostHandler) listPosts(w http.ResponseWriter, r *http.Request) {
	newsletterID := chi.URLParam(r, "newsletterId")

	posts, err := h.service.List(r.Context(), newsletterID)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
