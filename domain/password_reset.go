package domain

import "context"

// PasswordResetToken represents a reset token linked to a user.
type PasswordResetToken struct {
	Token     string
	UserID    string
	ExpiresAt int64 // unix timestamp for simplicity
}

// PasswordResetRepository defines storage operations for reset tokens.
type PasswordResetRepository interface {
	Create(ctx context.Context, t *PasswordResetToken) error
	Get(ctx context.Context, token string) (*PasswordResetToken, error)
	Delete(ctx context.Context, token string) error
}
