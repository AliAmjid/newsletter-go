package postgres

import (
	"context"
	"database/sql"

	"newsletter-go/domain"
)

// NewsletterRepository handles newsletter persistence in Postgres.
type NewsletterRepository struct {
	DB *sql.DB
}

func NewNewsletterRepository(db *sql.DB) *NewsletterRepository {
	return &NewsletterRepository{DB: db}
}

func (r *NewsletterRepository) ListByOwner(ctx context.Context, ownerID string) ([]*domain.Newsletter, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, name, description, owner_id, created_at FROM newsletter WHERE owner_id = $1`,
		ownerID,
	)
	if err != nil {
		return []*domain.Newsletter{}, err
	}
	defer rows.Close()

	var list []*domain.Newsletter
	for rows.Next() {
		var n domain.Newsletter
		if err := rows.Scan(&n.ID, &n.Name, &n.Description, &n.OwnerID, &n.CreatedAt); err != nil {
			return []*domain.Newsletter{}, err
		}
		list = append(list, &n)
	}
	if err := rows.Err(); err != nil {
		return []*domain.Newsletter{}, err
	}
	if list == nil {
		list = []*domain.Newsletter{}
	}
	return list, nil
}

func (r *NewsletterRepository) Create(ctx context.Context, n *domain.Newsletter) error {
	return r.DB.QueryRowContext(ctx,
		`INSERT INTO newsletter (name, description, owner_id) VALUES ($1, $2, $3) RETURNING id, created_at`,
		n.Name, n.Description, n.OwnerID,
	).Scan(&n.ID, &n.CreatedAt)
}

func (r *NewsletterRepository) GetByID(ctx context.Context, id string) (*domain.Newsletter, error) {
	var n domain.Newsletter
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, name, description, owner_id, created_at FROM newsletter WHERE id = $1`,
		id,
	).Scan(&n.ID, &n.Name, &n.Description, &n.OwnerID, &n.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NewsletterRepository) Update(ctx context.Context, n *domain.Newsletter) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE newsletter SET name = $1, description = $2 WHERE id = $3`,
		n.Name, n.Description, n.ID,
	)
	return err
}

func (r *NewsletterRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx,
		`DELETE FROM newsletter WHERE id = $1`,
		id,
	)
	return err
}
