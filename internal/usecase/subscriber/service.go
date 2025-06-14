package subscriber

import (
	"context"

	"github.com/google/uuid"

	"newsletter-go/domain"
)

// Service handles subscriber business logic.
type Service struct {
	repo domain.SubscriberRepository
}

func NewService(r domain.SubscriberRepository) *Service {
	return &Service{repo: r}
}

func (s *Service) List(ctx context.Context, newsletterID string) ([]*domain.Subscriber, error) {
	return s.repo.ListByNewsletter(ctx, newsletterID)
}

func (s *Service) Subscribe(ctx context.Context, newsletterID, email string) (*domain.Subscriber, error) {
	sub := &domain.Subscriber{NewsletterID: newsletterID, Email: email, Token: uuid.New().String()}
	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *Service) Confirm(ctx context.Context, token string) error {
	return s.repo.Confirm(ctx, token)
}

func (s *Service) Unsubscribe(ctx context.Context, token string) error {
	return s.repo.DeleteByToken(ctx, token)
}

func (s *Service) GetByToken(ctx context.Context, token string) (*domain.Subscriber, error) {
	return s.repo.GetByToken(ctx, token)
}
