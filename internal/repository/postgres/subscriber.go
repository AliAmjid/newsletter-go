package postgres

import (
	"context"
	"database/sql"

	"newsletter-go/domain"
)

// SubscriberRepository persists subscribers in Postgres.
type SubscriberRepository struct {
	DB *sql.DB
}

func NewSubscriberRepository(db *sql.DB) *SubscriberRepository {
	return &SubscriberRepository{DB: db}
}

func (r *SubscriberRepository) ListByNewsletter(ctx context.Context, newsletterID string) ([]*domain.Subscriber, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, newsletter_id, email, token, confirmed, created_at FROM subscriber WHERE newsletter_id = $1`,
		newsletterID,
	)
	if err != nil {
		return []*domain.Subscriber{}, err
	}
	defer rows.Close()

	var list []*domain.Subscriber
	for rows.Next() {
		var s domain.Subscriber
		if err := rows.Scan(&s.ID, &s.NewsletterID, &s.Email, &s.Token, &s.Confirmed, &s.CreatedAt); err != nil {
			return []*domain.Subscriber{}, err
		}
		list = append(list, &s)
	}
	if err := rows.Err(); err != nil {
		return []*domain.Subscriber{}, err
	}
	if list == nil {
		list = []*domain.Subscriber{}
	}
	return list, nil
}

func (r *SubscriberRepository) Create(ctx context.Context, s *domain.Subscriber) error {
	return r.DB.QueryRowContext(ctx,
		`INSERT INTO subscriber (newsletter_id, email, token) VALUES ($1, $2, $3) RETURNING id, confirmed, created_at`,
		s.NewsletterID, s.Email, s.Token,
	).Scan(&s.ID, &s.Confirmed, &s.CreatedAt)
}

func (r *SubscriberRepository) GetByToken(ctx context.Context, token string) (*domain.Subscriber, error) {
	var s domain.Subscriber
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, newsletter_id, email, token, confirmed, created_at FROM subscriber WHERE token = $1`,
		token,
	).Scan(&s.ID, &s.NewsletterID, &s.Email, &s.Token, &s.Confirmed, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubscriberRepository) Confirm(ctx context.Context, token string) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE subscriber SET confirmed = TRUE WHERE token = $1`,
		token,
	)
	return err
}

func (r *SubscriberRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.DB.ExecContext(ctx,
		`DELETE FROM subscriber WHERE token = $1`,
		token,
	)
	return err
}
