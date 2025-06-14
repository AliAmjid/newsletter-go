package domain

import (
	"context"
	"time"
)

// Subscriber represents an email subscribed to a newsletter.
type Subscriber struct {
	ID           string    `json:"id"`
	NewsletterID string    `json:"newsletterId"`
	Email        string    `json:"email"`
	Token        string    `json:"-"`
	Confirmed    bool      `json:"confirmed"`
	CreatedAt    time.Time `json:"createdAt"`
}

// SubscriberRepository defines persistence methods for subscribers.
type SubscriberRepository interface {
	ListByNewsletter(ctx context.Context, newsletterID string) ([]*Subscriber, error)
	Create(ctx context.Context, s *Subscriber) error
	GetByToken(ctx context.Context, token string) (*Subscriber, error)
	Confirm(ctx context.Context, token string) error
	DeleteByToken(ctx context.Context, token string) error
}
