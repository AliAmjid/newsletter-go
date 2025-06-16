package domain

import "context"

// PostDelivery links a post with a subscriber and tracks open state.
type PostDelivery struct {
	ID             string `json:"id"`
	PostID         string `json:"postId"`
	SubscriptionID string `json:"subscriptionId"`
	Opened         bool   `json:"opened"`
}

type PostDeliveryInfo struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Opened bool   `json:"opened"`
}

type PostDeliveryRepository interface {
	Create(ctx context.Context, postID, subscriptionID string) (*PostDelivery, error)
	MarkOpened(ctx context.Context, id string) error
	ListByPost(ctx context.Context, postID string) ([]*PostDeliveryInfo, error)
	ListByPostPaginated(ctx context.Context, postID, cursor string, limit int) ([]*PostDeliveryInfo, error)
}
