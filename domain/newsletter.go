package domain

import (
	"context"
	"time"
)

// Newsletter represents a newsletter owned by a user.
type Newsletter struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
}

// NewsletterRepository defines data access methods for newsletters.
type NewsletterRepository interface {
	ListByOwner(ctx context.Context, ownerID string) ([]*Newsletter, error)
	Create(ctx context.Context, n *Newsletter) error
	GetByID(ctx context.Context, id string) (*Newsletter, error)
	Update(ctx context.Context, n *Newsletter) error
	Delete(ctx context.Context, id string) error
	IsOwner(ctx context.Context, newsletterId, userId string) (bool, error)
}
