package router

import (
	"github.com/AliAmjid/newsletter-go/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Turn on logging
	r.Use(middleware.Logger)

	r.Route("/posts", func(r chi.Router) {
		r.Get("/{id}", handler.GetPost)
		r.Post(("/"), handler.CreatePost)
	})

	return r
}
