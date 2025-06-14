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

// SubscriptionRepository defines operations for newsletter subscriptions.
type SubscriptionRepository interface {
	Create(ctx context.Context, s *Subscription) error
	Confirm(ctx context.Context, token string) (*Subscription, error)
	DeleteByToken(ctx context.Context, token string) error
	ListByNewsletter(ctx context.Context, newsletterID string) ([]*Subscription, error)
	GetByNewsletterEmail(ctx context.Context, newsletterID, email string) (*Subscription, error)
	UpdateToken(ctx context.Context, id, token string) (*Subscription, error)
}
