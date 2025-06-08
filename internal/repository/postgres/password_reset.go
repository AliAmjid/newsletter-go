package postgres

import (
	"context"
	"database/sql"
	"time"

	"newsletter-go/domain"
)

// PasswordResetRepository stores reset tokens in Postgres.
type PasswordResetRepository struct {
	DB *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{DB: db}
}

func (r *PasswordResetRepository) Create(ctx context.Context, t *domain.PasswordResetToken) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO password_reset_tokens (token, user_id, expires_at) VALUES ($1, $2, to_timestamp($3))`,
		t.Token, t.UserID, t.ExpiresAt,
	)
	return err
}

func (r *PasswordResetRepository) Get(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
	var t domain.PasswordResetToken
	var expires time.Time
	err := r.DB.QueryRowContext(ctx,
		`SELECT token, user_id, expires_at FROM password_reset_tokens WHERE token = $1`,
		token,
	).Scan(&t.Token, &t.UserID, &expires)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t.ExpiresAt = expires.Unix()
	return &t, nil
}

func (r *PasswordResetRepository) Delete(ctx context.Context, token string) error {
	_, err := r.DB.ExecContext(ctx,
		`DELETE FROM password_reset_tokens WHERE token = $1`,
		token,
	)
	return err
}
