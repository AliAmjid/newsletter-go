package subscriber

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetSubscribers(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 4 || parts[1] != "newsletters" || parts[3] != "subscribers" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	newsletterID := parts[2]

	rows, err := h.db.Query(`
		SELECT id, email, subscribed_at
		FROM subscribers
		WHERE newsletter_id = $1
	`, newsletterID)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var subscribers []Subscriber
	for rows.Next() {
		var s Subscriber
		if err := rows.Scan(&s.ID, &s.Email, &s.SubscribedAt); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		subscribers = append(subscribers, s)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Row error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscribers)
}
