package repository

import (
	"context"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct{ db *pgxpool.Pool }

func NewEventRepository(db *pgxpool.Pool) *EventRepository { return &EventRepository{db: db} }

func (r *EventRepository) Create(ctx context.Context, e model.Event) (model.Event, error) {
	_, err := r.db.Exec(ctx, `INSERT INTO events (id, subscription_id, type, payload) VALUES ($1, $2, $3, $4)`, e.ID, e.SubscriptionID, e.Type, e.Payload)
	return e, err
}

func (r *EventRepository) ListBySubscriptionID(ctx context.Context, subscriptionID string) ([]model.Event, error) {
	rows, err := r.db.Query(ctx, `SELECT id, subscription_id, type, payload, created_at FROM events WHERE subscription_id = $1 ORDER BY created_at ASC`, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.Event
	for rows.Next() {
		var e model.Event
		if err := rows.Scan(&e.ID, &e.SubscriptionID, &e.Type, &e.Payload, &e.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}
