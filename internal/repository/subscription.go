package repository

import (
	"context"
	"errors"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct{ db *pgxpool.Pool }

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, s model.Subscription) (model.Subscription, error) {
	_, err := r.db.Exec(ctx, `INSERT INTO subscriptions (id, customer_id, plan_code, amount_cents, currency, status) VALUES ($1, $2, $3, $4, $5, $6)`, s.ID, s.CustomerID, s.PlanCode, s.AmountCents, s.Currency, s.Status)
	return s, err
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (model.Subscription, error) {
	var s model.Subscription
	err := r.db.QueryRow(ctx, `SELECT id, customer_id, plan_code, amount_cents, currency, status, created_at FROM subscriptions WHERE id = $1`, id).Scan(&s.ID, &s.CustomerID, &s.PlanCode, &s.AmountCents, &s.Currency, &s.Status, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Subscription{}, shared.ErrNotFound
		}
		return model.Subscription{}, err
	}
	return s, nil
}

func (r *SubscriptionRepository) UpdateStatus(ctx context.Context, id string, status model.SubscriptionStatus) error {
	ct, err := r.db.Exec(ctx, `UPDATE subscriptions SET status = $2 WHERE id = $1`, id, status)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return shared.ErrNotFound
	}
	return nil
}
