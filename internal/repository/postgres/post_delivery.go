package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"newsletter-go/domain"
)

type PostDeliveryRepository struct {
	DB *sql.DB
}

func NewPostDeliveryRepository(db *sql.DB) *PostDeliveryRepository {
	return &PostDeliveryRepository{DB: db}
}

func (r *PostDeliveryRepository) Create(ctx context.Context, postID, subscriptionID string) (*domain.PostDelivery, error) {
	var d domain.PostDelivery
	err := r.DB.QueryRowContext(ctx,
		`INSERT INTO post_delivery (post_id, subscription_id) VALUES ($1,$2) RETURNING id, opened`,
		postID, subscriptionID,
	).Scan(&d.ID, &d.Opened)
	if err != nil {
		return nil, err
	}
	d.PostID = postID
	d.SubscriptionID = subscriptionID
	return &d, nil
}

func (r *PostDeliveryRepository) MarkOpened(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `UPDATE post_delivery SET opened = TRUE WHERE id = $1 AND opened = FALSE`, id)
	return err
}

func (r *PostDeliveryRepository) ListByPost(ctx context.Context, postID string) ([]*domain.PostDeliveryInfo, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT d.id, s.email, d.opened FROM post_delivery d JOIN subscription s ON s.id = d.subscription_id WHERE d.post_id = $1`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var infos []*domain.PostDeliveryInfo
	for rows.Next() {
		var info domain.PostDeliveryInfo
		if err := rows.Scan(&info.ID, &info.Email, &info.Opened); err != nil {
			return nil, err
		}
		infos = append(infos, &info)
	}
	if infos == nil {
		infos = []*domain.PostDeliveryInfo{}
	}
	return infos, rows.Err()
}

func (r *PostDeliveryRepository) ListByPostPaginated(ctx context.Context, postID, cursor string, limit int) ([]*domain.PostDeliveryInfo, error) {
	args := []interface{}{postID}
	query := `SELECT d.id, s.email, d.opened FROM post_delivery d JOIN subscription s ON s.id = d.subscription_id WHERE d.post_id = $1`
	idx := 2
	if cursor != "" {
		query += fmt.Sprintf(" AND d.id > $%d", idx)
		args = append(args, cursor)
		idx++
	}
	query += fmt.Sprintf(" ORDER BY d.id LIMIT $%d", idx)
	args = append(args, limit)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var infos []*domain.PostDeliveryInfo
	for rows.Next() {
		var info domain.PostDeliveryInfo
		if err := rows.Scan(&info.ID, &info.Email, &info.Opened); err != nil {
			return nil, err
		}
		infos = append(infos, &info)
	}
	if infos == nil {
		infos = []*domain.PostDeliveryInfo{}
	}
	return infos, rows.Err()
}
