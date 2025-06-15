package domain

import "context"

// PostDelivery links a post with a subscriber and tracks open state.
type PostDelivery struct {
	ID             string
	PostID         string
	SubscriptionID string
	Opened         bool
}

type PostDeliveryInfo struct {
	ID     string
	Email  string
	Opened bool
}

type PostDeliveryRepository interface {
	Create(ctx context.Context, postID, subscriptionID string) (*PostDelivery, error)
	MarkOpened(ctx context.Context, id string) error
	ListByPost(ctx context.Context, postID string) ([]*PostDeliveryInfo, error)
	ListByPostPaginated(ctx context.Context, postID, cursor string, limit int) ([]*PostDeliveryInfo, error)
}
