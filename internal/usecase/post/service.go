package post

import (
	"context"
	"errors"
	"time"

	"newsletter-go/domain"
)

// Service handles post business logic.
type Service struct {
	repo           domain.PostRepository
	newsletterRepo domain.NewsletterRepository
}

var (
	ErrNotFound         = errors.New("post not found")
	ErrNotOwner         = errors.New("not the owner")
	ErrAlreadyPublished = errors.New("post already published")
)

// NewService creates a Service instance.
func NewService(r domain.PostRepository, nR domain.NewsletterRepository) *Service {
	return &Service{
		repo:           r,
		newsletterRepo: nR,
	}
}

// Save persists a post via the repository.
func (s *Service) Create(ctx context.Context, userID string, p *domain.Post) error {
	isOwner, err := s.newsletterRepo.IsOwner(ctx, p.NewsletterId, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrNotOwner
	}
	return s.repo.Create(ctx, p)
}

func (s *Service) Publish(ctx context.Context, userID, id string) (*domain.Post, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, ErrNotFound
	}

	isOwner, err := s.newsletterRepo.IsOwner(ctx, post.NewsletterId, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, ErrNotOwner
	}

	if post.PublishedAt != nil {
		return nil, ErrAlreadyPublished
	}

	return s.repo.Publish(ctx, id)
}

func (s *Service) IsNewsletterOwner(ctx context.Context, newsletterId, userId string) (bool, error) {
	return s.newsletterRepo.IsOwner(ctx, newsletterId, userId)
}

func (s *Service) List(ctx context.Context, newsletterId, cursor string, limit int) ([]*domain.Post, string, error) {
	posts, err := s.repo.ListByNewsletter(ctx, newsletterId, cursor, limit)
	if err != nil {
		return nil, "", err
	}
	next := ""
	if len(posts) == limit {
		next = posts[len(posts)-1].PublishedAt.Format(time.RFC3339)
	}
	return posts, next, nil
}
