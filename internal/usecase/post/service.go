package post

import (
	"context"
	"errors"
	"time"

	"newsletter-go/domain"
	"newsletter-go/internal/mailer"
)

// Service handles post business logic.
type Service struct {
	repo           domain.PostRepository
	newsletterRepo domain.NewsletterRepository
	subRepo        domain.SubscriptionRepository
	deliveryRepo   domain.PostDeliveryRepository
	mailer         *mailer.Service
}

var (
	ErrNotFound         = errors.New("post not found")
	ErrNotOwner         = errors.New("not the owner")
	ErrAlreadyPublished = errors.New("post already published")
)

// NewService creates a Service instance.
func NewService(r domain.PostRepository, nR domain.NewsletterRepository, subRepo domain.SubscriptionRepository, dRepo domain.PostDeliveryRepository, m *mailer.Service) *Service {
	return &Service{
		repo:           r,
		newsletterRepo: nR,
		subRepo:        subRepo,
		deliveryRepo:   dRepo,
		mailer:         m,
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

	post, err = s.repo.Publish(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.mailer != nil && s.subRepo != nil {
		go s.notifySubscribers(post)
	}
	return post, nil
}

func (s *Service) notifySubscribers(p *domain.Post) {
	subs, err := s.subRepo.ListByNewsletterAll(context.Background(), p.NewsletterId)
	if err != nil {
		return
	}
	for _, sub := range subs {
		if sub.ConfirmedAt == nil {
			continue
		}
		delivery, err := s.deliveryRepo.Create(context.Background(), p.ID, sub.ID)
		if err != nil {
			continue
		}
		email := sub.Email
		token := sub.Token
		deliveryID := delivery.ID
		go func() {
			_ = s.mailer.SendPostEmail(email, token, p, deliveryID)
		}()
	}
}

func (s *Service) MarkOpened(ctx context.Context, deliveryID string) error {
	return s.deliveryRepo.MarkOpened(ctx, deliveryID)
}

type Metrics struct {
	TotalSend   int
	TotalOpened int
	Deliveries  []*domain.PostDeliveryInfo
}

func (s *Service) GetWithMetrics(ctx context.Context, userID, postID string) (*domain.Post, *Metrics, error) {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, nil, err
	}
	if post == nil {
		return nil, nil, ErrNotFound
	}
	isOwner, err := s.newsletterRepo.IsOwner(ctx, post.NewsletterId, userID)
	if err != nil {
		return nil, nil, err
	}
	if !isOwner {
		return nil, nil, ErrNotOwner
	}
	deliveries, err := s.deliveryRepo.ListByPost(ctx, postID)
	if err != nil {
		return nil, nil, err
	}
	m := &Metrics{Deliveries: deliveries, TotalSend: len(deliveries)}
	for _, d := range deliveries {
		if d.Opened {
			m.TotalOpened++
		}
	}
	return post, m, nil
}

func (s *Service) IsNewsletterOwner(ctx context.Context, newsletterId, userId string) (bool, error) {
	return s.newsletterRepo.IsOwner(ctx, newsletterId, userId)
}

func (s *Service) List(ctx context.Context, newsletterId, cursor string, limit int, search string) ([]*domain.Post, string, error) {
	posts, err := s.repo.ListByNewsletter(ctx, newsletterId, cursor, limit, search)
	if err != nil {
		return nil, "", err
	}
	next := ""
	if len(posts) == limit {
		next = posts[len(posts)-1].PublishedAt.Format(time.RFC3339)
	}
	return posts, next, nil
}

func (s *Service) ListDeliveries(ctx context.Context, userID, postID, cursor string, limit int) ([]*domain.PostDeliveryInfo, string, error) {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, "", err
	}
	if post == nil {
		return nil, "", ErrNotFound
	}
	ok, err := s.newsletterRepo.IsOwner(ctx, post.NewsletterId, userID)
	if err != nil {
		return nil, "", err
	}
	if !ok {
		return nil, "", ErrNotOwner
	}

	infos, err := s.deliveryRepo.ListByPostPaginated(ctx, postID, cursor, limit)
	if err != nil {
		return nil, "", err
	}

	next := ""
	if len(infos) == limit {
		next = infos[len(infos)-1].ID
	}
	return infos, next, nil
}
