package postgres

import (
	"context"
	"database/sql"

	"newsletter-go/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	return r.DB.QueryRowContext(ctx,
		`INSERT INTO "user" (email, password_hash, firebase_uid) VALUES ($1, $2, $3) RETURNING id`,
		u.Email, u.PasswordHash, u.FirebaseUID,
	).Scan(&u.ID)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, email, password_hash, firebase_uid FROM "user" WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirebaseUID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, email, password_hash, firebase_uid FROM "user" WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirebaseUID)
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

func (r *UserRepository) GetByFirebaseID(ctx context.Context, fid string) (*domain.User, error) {
	var u domain.User
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, email, password_hash, firebase_uid FROM "user" WHERE firebase_uid = $1`,
		fid,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirebaseUID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
