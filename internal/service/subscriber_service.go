package service

import (
	"context"
	"newsletter-go/domain"
)

type SubscriberService struct {
	repo domain.UserRepository
}

func NewSubscriberService(repo domain.UserRepository) *SubscriberService {
	return &SubscriberService{repo: repo}
}

func (s *SubscriberService) Create(ctx context.Context, u *domain.User) error {
	return s.repo.Create(ctx, u)
}

func (s *SubscriberService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}
