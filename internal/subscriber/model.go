// internal/subscriber/model.go

package subscriber

import "time"

type Subscriber struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	SubscribedAt time.Time `json:"subscribedAt"`
}
