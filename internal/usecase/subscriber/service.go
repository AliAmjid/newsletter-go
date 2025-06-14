package subscriber

import (
	"context"
	"newsletter-go/domain"
)

type Service struct {
	repo domain.UserRepository
}

func NewService(repo domain.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, u *domain.User) error {
	return s.repo.Create(ctx, u)
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}
