package postgres

import (
	"context"
	"database/sql"

	"github.com/AliAmjid/newsletter-go/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	return r.DB.QueryRowContext(ctx,
		`INSERT INTO "user" (email, password_hash) VALUES ($1, $2) RETURNING id`,
		u.Email, u.PasswordHash,
	).Scan(&u.ID)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, email, password_hash FROM "user" WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id string, hash string) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE "user" SET password_hash = $1 WHERE id = $2`,
		hash, id,
	)
	return err
}
