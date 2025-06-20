package http

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

var pixelData []byte

func init() {
	b, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMB/6X9ZPkAAAAASUVORK5CYII=")
	pixelData = b
}

type postCreateRequest struct {
	Title              string `json:"title" validate:"required"`
	Content            string `json:"content"`
	PublishImmediately bool   `json:"publishImmediately"`
}

type PaginatedPostResponse struct {
	Posts      []*domain.Post `json:"posts"`
	NextCursor string         `json:"nextCursor"`
}

type SinglePostResponse struct {
	Post        *domain.Post `json:"post"`
	TotalSend   int          `json:"totalSend"`
	TotalOpened int          `json:"totalOpened"`
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
		r.Post("/{postId}/publish", h.publishPost)
		r.Get("/{postId}", h.getPost)
	})
	r.Get("/posts/{postId}/deliveries", h.listDeliveries)
	r.Get("/post-deliveries/{deliveryId}/pixel", h.pixel)
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

	var req postCreateRequest
	if !bindAndValidate(w, r, &req, h.validate) {
		return
	}

	p := &domain.Post{
		NewsletterId: newsletterId,
		Title:        req.Title,
		Content:      req.Content,
	}
	if err := h.service.Create(r.Context(), user.ID, p); err != nil {
		if errors.Is(err, postusecase.ErrNotOwner) {
			respondWithError(w, http.StatusForbidden, err.Error())
			return
		}
		fmt.Printf("%+v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to save post")
		return
	}
	if req.PublishImmediately {
		pub, err := h.service.Publish(r.Context(), user.ID, p.ID)
		if err != nil {
			fmt.Printf("%+v\n", err)
		} else {
			p = pub
		}
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

	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	pubStr := r.URL.Query().Get("published")
	var pubFilter *bool
	switch pubStr {
	case "1":
		b := true
		pubFilter = &b
	case "0":
		b := false
		pubFilter = &b
	case "", "null":
	default:
		respondWithError(w, http.StatusBadRequest, "invalid published value")
		return
	}
	limit := 20
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			limit = v
		}
	}

	posts, next, err := h.service.List(r.Context(), newsletterId, cursor, limit, search, pubFilter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list posts")
		return
	}

	result := PaginatedPostResponse{
		Posts:      posts,
		NextCursor: next,
	}
	respondWithJSON(w, http.StatusOK, result)
}

func (h *PostHandler) publishPost(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")

	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if ok, err := h.users.IsAllowedTo(r, "update", "post"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	p, err := h.service.Publish(r.Context(), user.ID, postId)
	if err != nil {
		if errors.Is(err, postusecase.ErrNotOwner) {
			respondWithError(w, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, postusecase.ErrAlreadyPublished) {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, postusecase.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to publish post")
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (h *PostHandler) pixel(w http.ResponseWriter, r *http.Request) {
	deliveryId := chi.URLParam(r, "deliveryId")
	if deliveryId != "" {
		_ = h.service.MarkOpened(r.Context(), deliveryId)
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(pixelData)
}

func (h *PostHandler) getPost(w http.ResponseWriter, r *http.Request) {
	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if ok, err := h.users.IsAllowedTo(r, "read", "post"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	postId := chi.URLParam(r, "postId")
	p, m, err := h.service.GetWithMetrics(r.Context(), user.ID, postId)
	if err != nil {
		if err == postusecase.ErrNotOwner {
			respondWithError(w, http.StatusForbidden, "You are not the owner of this newsletter")
			return
		}
		if err == postusecase.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch post")
		return
	}
	result := SinglePostResponse{
		Post:        p,
		TotalSend:   m.TotalSend,
		TotalOpened: m.TotalOpened,
	}
	respondWithJSON(w, http.StatusOK, result)
}

func (h *PostHandler) listDeliveries(w http.ResponseWriter, r *http.Request) {
	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if ok, err := h.users.IsAllowedTo(r, "read", "post"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	postID := chi.URLParam(r, "postId")
	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			limit = v
		}
	}

	deliveries, next, err := h.service.ListDeliveries(r.Context(), user.ID, postID, cursor, limit)
	if err != nil {
		if errors.Is(err, postusecase.ErrNotOwner) {
			respondWithError(w, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, postusecase.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to list deliveries")
		return
	}

	type deliveriesResponse struct {
		Deliveries []*domain.PostDeliveryInfo `json:"deliveries"`
		NextCursor string                     `json:"nextCursor"`
	}

	respondWithJSON(w, http.StatusOK, deliveriesResponse{Deliveries: deliveries, NextCursor: next})
}
