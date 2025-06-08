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

}

func (h *HelloHandler) sayHello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello!"))
}
