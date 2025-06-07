package service

import (
	"github.com/AliAmjid/newsletter-go/internal/db"
	"github.com/AliAmjid/newsletter-go/internal/model"
)

func SavePost(post model.Post) error {
	_, err := db.DB.Exec(
		"INSERT INTO posts (title, content) VALUES ($1, $2)",
		post.Title, post.Content,
	)

	return err
}
