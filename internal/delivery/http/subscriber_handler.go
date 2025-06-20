package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	subscriberusecase "newsletter-go/internal/usecase/subscriber"
	userusecase "newsletter-go/internal/usecase/user"
)

// SubscriberHandler handles subscription endpoints.
type SubscriberHandler struct {
	service  *subscriberusecase.Service
	users    *userusecase.Service
	validate *validator.Validate
}

func NewSubscriberHandler(r chi.Router, s *subscriberusecase.Service, u *userusecase.Service) {
	h := &SubscriberHandler{service: s, users: u, validate: validator.New()}

	r.Route("/subscriptions/{newsletterId}", func(r chi.Router) {
		r.Post("/subscribe", h.subscribe)
		r.Get("/subscribers", h.listSubscribers)
	})

	r.Get("/subscriptions/confirm", h.confirm)
	r.Get("/subscriptions/unsubscribe", h.unsubscribe)
}

type subscribeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (h *SubscriberHandler) subscribe(w http.ResponseWriter, r *http.Request) {
	newsletterID := chi.URLParam(r, "newsletterId")
	var req subscribeRequest
	if !bindAndValidate(w, r, &req, h.validate) {
		return
	}
	if _, err := h.service.Subscribe(r.Context(), newsletterID, req.Email); err != nil {
		if err == subscriberusecase.ErrTooFrequent {
			respondWithError(w, http.StatusTooManyRequests, err.Error())
			return
		}
		if err == subscriberusecase.ErrAlreadySubscribed {
			respondWithError(w, http.StatusConflict, err.Error())
			return
		}

		respondWithError(w, http.StatusInternalServerError, "failed to subscribe")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *SubscriberHandler) confirm(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respondWithError(w, http.StatusBadRequest, "token required")
		return
	}
	sub, err := h.service.Confirm(r.Context(), token)
	if err != nil || sub == nil {
		respondWithError(w, http.StatusBadRequest, "invalid or expired token")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *SubscriberHandler) unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respondWithError(w, http.StatusBadRequest, "token required")
		return
	}
	if err := h.service.Unsubscribe(r.Context(), token); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid token")
		return
	}
	w.WriteHeader(http.StatusOK)
}

type SubscriberResponse struct {
	Email       string     `json:"email"`
	ConfirmedAt *time.Time `json:"confirmedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt,omitempty"`
}

type PaginatedSubscriberResponse struct {
	Subscribers []SubscriberResponse `json:"subscribers"`
	NextCursor  string               `json:"nextCursor"`
}

func (h *SubscriberHandler) listSubscribers(w http.ResponseWriter, r *http.Request) {
	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if ok, err := h.users.IsAllowedTo(r, "read", "subscriber"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}
	newsletterID := chi.URLParam(r, "newsletterId")
	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	limit := 20
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			limit = v
		}
	}

	subs, next, err := h.service.List(r.Context(), newsletterID, cursor, limit, search)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list subscribers")
		return
	}

	respSubs := make([]SubscriberResponse, 0, len(subs))
	for _, s := range subs {
		respSubs = append(respSubs, SubscriberResponse{
			Email:       s.Email,
			ConfirmedAt: s.ConfirmedAt,
			CreatedAt:   s.CreatedAt,
		})
	}

	result := PaginatedSubscriberResponse{
		Subscribers: respSubs,
		NextCursor:  next,
	}

	respondWithJSON(w, http.StatusOK, result)
}
