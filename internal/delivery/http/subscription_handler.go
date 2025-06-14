package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	newslettercase "newsletter-go/internal/usecase/newsletter"
	subscribercase "newsletter-go/internal/usecase/subscriber"
	userusecase "newsletter-go/internal/usecase/user"
)

// SubscriptionHandler manages newsletter subscriptions.
type SubscriptionHandler struct {
	service     *subscribercase.Service
	newsletters *newslettercase.Service
	users       *userusecase.Service
	validate    *validator.Validate
}

func NewSubscriptionHandler(r chi.Router, s *subscribercase.Service, n *newslettercase.Service, u *userusecase.Service) {
	h := &SubscriptionHandler{service: s, newsletters: n, users: u, validate: validator.New()}

	r.Route("/newsletters", func(r chi.Router) {
		r.Post("/{newsletterId}/subscribe", h.subscribe)
		r.Get("/{newsletterId}/subscribers", h.listSubscribers)
	})
	r.Get("/subscriptions/confirm", h.confirm)
	r.Get("/subscriptions/unsubscribe", h.unsubscribe)
}

type subscribeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (h *SubscriptionHandler) subscribe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "newsletterId")
	var req subscribeRequest
	if !bindAndValidate(w, r, &req, h.validate) {
		return
	}
	_, err := h.service.Subscribe(r.Context(), id, req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to subscribe")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *SubscriptionHandler) confirm(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respondWithError(w, http.StatusBadRequest, "missing token")
		return
	}
	sub, err := h.service.GetByToken(r.Context(), token)
	if err != nil || sub == nil {
		respondWithError(w, http.StatusBadRequest, "invalid token")
		return
	}
	if err := h.service.Confirm(r.Context(), token); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to confirm")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *SubscriptionHandler) unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respondWithError(w, http.StatusBadRequest, "missing token")
		return
	}
	if err := h.service.Unsubscribe(r.Context(), token); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to unsubscribe")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *SubscriptionHandler) listSubscribers(w http.ResponseWriter, r *http.Request) {
	user, err := h.users.IsLoggedIn(r)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if ok, err := h.users.IsAllowedTo(r, "read", "subscribers"); err != nil || !ok {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}
	id := chi.URLParam(r, "newsletterId")
	nl, err := h.newsletters.GetByID(r.Context(), id)
	if err != nil || nl == nil {
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}
	if nl.OwnerID != user.ID {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	list, err := h.service.List(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list subscribers")
		return
	}
	respondWithJSON(w, http.StatusOK, list)
}
