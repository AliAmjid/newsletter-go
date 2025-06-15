package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"newsletter-go/domain"
)

type SubscriptionRepository struct {
	DB *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, s *domain.Subscription) error {
	return r.DB.QueryRowContext(ctx,
		`INSERT INTO subscription (newsletter_id, email, token) VALUES ($1,$2,$3) RETURNING id, created_at`,
		s.NewsletterID, s.Email, s.Token,
	).Scan(&s.ID, &s.CreatedAt)
}

func (r *SubscriptionRepository) Confirm(ctx context.Context, token string) (*domain.Subscription, error) {
	var s domain.Subscription
	err := r.DB.QueryRowContext(ctx,
		`UPDATE subscription SET confirmed_at = NOW() WHERE token = $1 AND confirmed_at IS NULL RETURNING id, newsletter_id, email, token, confirmed_at, created_at`,
		token,
	).Scan(&s.ID, &s.NewsletterID, &s.Email, &s.Token, &s.ConfirmedAt, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubscriptionRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM subscription WHERE token = $1`, token)
	return err
}

func (r *SubscriptionRepository) ListByNewsletter(ctx context.Context, newsletterID, cursor string, limit int, search string) ([]*domain.Subscription, error) {
	args := []interface{}{newsletterID}
	query := `SELECT id, newsletter_id, email, token, confirmed_at, created_at
                        FROM subscription
                        WHERE newsletter_id = $1 AND confirmed_at IS NOT NULL`
	idx := 2
	if cursor != "" {
		query += fmt.Sprintf(" AND created_at < $%d", idx)
		args = append(args, cursor)
		idx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND email ILIKE '%%' || $%d || '%%'", idx)
		args = append(args, search)
		idx++
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", idx)
	args = append(args, limit)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []*domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.NewsletterID, &s.Email, &s.Token, &s.ConfirmedAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, &s)
	}
	if subs == nil {
		subs = []*domain.Subscription{}
	}
	return subs, rows.Err()
}

func (r *SubscriptionRepository) ListByNewsletterAll(ctx context.Context, newsletterID string) ([]*domain.Subscription, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, newsletter_id, email, token, confirmed_at, created_at
                 FROM subscription
                 WHERE newsletter_id = $1 AND confirmed_at IS NOT NULL`,
		newsletterID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []*domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.NewsletterID, &s.Email, &s.Token, &s.ConfirmedAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, &s)
	}
	if subs == nil {
		subs = []*domain.Subscription{}
	}
	return subs, rows.Err()
}

func (r *SubscriptionRepository) GetByNewsletterEmail(ctx context.Context, newsletterID, email string) (*domain.Subscription, error) {
	var s domain.Subscription
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, newsletter_id, email, token, confirmed_at, created_at FROM subscription WHERE newsletter_id = $1 AND email = $2`,
		newsletterID, email,
	).Scan(&s.ID, &s.NewsletterID, &s.Email, &s.Token, &s.ConfirmedAt, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubscriptionRepository) UpdateToken(ctx context.Context, id, token string) (*domain.Subscription, error) {
	var s domain.Subscription
	err := r.DB.QueryRowContext(ctx,
		`UPDATE subscription SET token = $1, created_at = NOW() WHERE id = $2 RETURNING id, newsletter_id, email, token, confirmed_at, created_at`,
		token, id,
	).Scan(&s.ID, &s.NewsletterID, &s.Email, &s.Token, &s.ConfirmedAt, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
