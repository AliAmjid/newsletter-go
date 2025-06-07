package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("TODO post with ID " + id))
}
