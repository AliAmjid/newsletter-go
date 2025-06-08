package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	usersvc "newsletter-go/internal/usecase/user"
)

type HelloHandler struct {
	users *usersvc.Service
}

func NewHelloHandler(r chi.Router, us *usersvc.Service) {
	h := &HelloHandler{users: us}

	r.Get("/", h.sayHello)
	r.Get("/whoami", h.whoAmI)
}

func (h *HelloHandler) sayHello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello!"))
}

func (h *HelloHandler) whoAmI(w http.ResponseWriter, r *http.Request) {
	u, err := h.users.IsLoggedIn(r)
	if err != nil || u == nil {
		respondWithError(w, http.StatusUnauthorized, "not logged in")
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}
