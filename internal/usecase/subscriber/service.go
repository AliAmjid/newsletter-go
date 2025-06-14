package subscriber

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"newsletter-go/domain"
	"newsletter-go/internal/mailer"
)

type Service struct {
	repo   domain.SubscriptionRepository
	mailer *mailer.Service
}

var ErrTooFrequent = errors.New("confirmation email sent recently")
var ErrAlreadySubscribed = errors.New("already subscribed")

func NewService(r domain.SubscriptionRepository, m *mailer.Service) *Service {
	return &Service{repo: r, mailer: m}
}

func (s *Service) Subscribe(ctx context.Context, newsletterID, email string) (string, error) {
	existing, err := s.repo.GetByNewsletterEmail(ctx, newsletterID, email)
	if err != nil {
		return "", err
	}

	if existing.ConfirmedAt != nil {
		return "", ErrAlreadySubscribed
	}

	token := uuid.New().String()

	if existing != nil {
		if time.Since(existing.CreatedAt) < time.Minute {
			return "", ErrTooFrequent
		}
		if _, err := s.repo.UpdateToken(ctx, existing.ID, token); err != nil {
			return "", err
		}
	} else {
		sub := &domain.Subscription{NewsletterID: newsletterID, Email: email, Token: token}
		if err := s.repo.Create(ctx, sub); err != nil {
			return "", err
		}
	}

	if s.mailer != nil {
		_ = s.mailer.SendSubscriptionConfirmEmail(email, token)
	}
	return token, nil
}

func (s *Service) Confirm(ctx context.Context, token string) (*domain.Subscription, error) {
	return s.repo.Confirm(ctx, token)
}

func (s *Service) Unsubscribe(ctx context.Context, token string) error {
	return s.repo.DeleteByToken(ctx, token)
}

func (s *Service) List(ctx context.Context, newsletterID string) ([]*domain.Subscription, error) {
	return s.repo.ListByNewsletter(ctx, newsletterID)
}
