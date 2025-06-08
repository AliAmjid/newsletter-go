package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AliAmjid/newsletter-go/internal/model"
	"github.com/AliAmjid/newsletter-go/internal/service"
	"github.com/go-chi/chi/v5"
)

func GetPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("TODO post with ID " + id))
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post model.Post

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := service.SavePost(post); err != nil {
		http.Error(w, "Failed to save post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Post recieved"))
}
