package domain

import (
	"context"
	"time"
)

// Subscription represents an email subscribed to a newsletter.
type Subscription struct {
	ID           string
	NewsletterID string
	Email        string
	Token        string
	ConfirmedAt  *time.Time
	CreatedAt    time.Time
}

type SubscriptionRepository interface {
	Create(ctx context.Context, s *Subscription) error
	Confirm(ctx context.Context, token string) (*Subscription, error)
	DeleteByToken(ctx context.Context, token string) error
	// ListByNewsletter returns confirmed subscriptions for the given newsletter.
	// Results are ordered by creation date descending. If cursor is provided,
	// only subscriptions created before the cursor will be returned. The
	// number of results is limited by limit.
	ListByNewsletter(ctx context.Context, newsletterID, cursor string, limit int, search string) ([]*Subscription, error)
	ListByNewsletterAll(ctx context.Context, newsletterID string) ([]*Subscription, error)
	GetByNewsletterEmail(ctx context.Context, newsletterID, email string) (*Subscription, error)
	UpdateToken(ctx context.Context, id, token string) (*Subscription, error)
}
