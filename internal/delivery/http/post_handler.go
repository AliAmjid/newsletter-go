package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"newsletter-go/domain"
	postusecase "newsletter-go/internal/usecase/post"
	userusecase "newsletter-go/internal/usecase/user"
)

type PostHandler struct {
	service  *postusecase.Service
	users    *userusecase.Service
	validate *validator.Validate
}

type postCreateRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
}

func NewPostHandler(r chi.Router, s *postusecase.Service, u *userusecase.Service) {
	h := &PostHandler{
		service:  s,
		users:    u,
		validate: validator.New(),
	}

	r.Route("/newsletters/{newsletterId}/posts", func(r chi.Router) {
		r.Get("/", h.listPosts)
		r.Post("/", h.createPost)
	})
}

func (h *PostHandler) createPost(w http.ResponseWriter, r *http.Request) {
	newsletterId := chi.URLParam(r, "newsletterId")

	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if ok, err := h.users.IsAllowedTo(r, "create", "post"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	if isOwner, err := h.service.IsNewsletterOwner(r.Context(), newsletterId, user.ID); err != nil || !isOwner {
		respondWithError(w, http.StatusForbidden, "You are not the owner of this newsletter")
		return
	}

	var req postCreateRequest
	if !bindAndValidate(w, r, &req, h.validate) {
		return
	}

	p := &domain.Post{
		NewsletterId: newsletterId,
		Title:        req.Title,
		Content:      req.Content,
	}
	if err := h.service.Save(r.Context(), p); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save post")
		return
	}
	respondWithJSON(w, http.StatusCreated, p)
}

func (h *PostHandler) listPosts(w http.ResponseWriter, r *http.Request) {

	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if ok, err := h.users.IsAllowedTo(r, "read", "post"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	newsletterId := chi.URLParam(r, "newsletterId")
	if newsletterId == "" {
		respondWithError(w, http.StatusBadRequest, "newsletterId is required")
		return
	}

	posts, err := h.service.ListPostsByNewsletter(r.Context(), newsletterId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list posts")
		return
	}

	respondWithJSON(w, http.StatusOK, posts)

}
