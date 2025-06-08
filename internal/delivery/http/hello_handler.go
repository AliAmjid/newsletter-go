package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HelloHandler struct{}

func NewHelloHandler(r chi.Router) {
	h := &HelloHandler{}

	r.Get("/", h.sayHello)
}

func (h *HelloHandler) sayHello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello!"))
}
