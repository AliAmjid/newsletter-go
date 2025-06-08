package domain

import "context"

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdatePassword(ctx context.Context, id string, hash string) error
}
