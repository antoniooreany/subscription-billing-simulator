package repository

import (
	"context"
	"errors"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct{ db *pgxpool.Pool }

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository { return &PaymentRepository{db: db} }

func (r *PaymentRepository) Create(ctx context.Context, p model.Payment) (model.Payment, error) {
	_, err := r.db.Exec(ctx, `INSERT INTO payments (id, subscription_id, amount_cents, currency, status, reason, idempotency_key) VALUES ($1, $2, $3, $4, $5, $6, $7)`, p.ID, p.SubscriptionID, p.AmountCents, p.Currency, p.Status, p.Reason, p.IdempotencyKey)
	return p, err
}

func (r *PaymentRepository) GetByIdempotencyKey(ctx context.Context, key string) (model.Payment, error) {
	var p model.Payment
	err := r.db.QueryRow(ctx, `SELECT id, subscription_id, amount_cents, currency, status, reason, idempotency_key, created_at FROM payments WHERE idempotency_key = $1`, key).Scan(&p.ID, &p.SubscriptionID, &p.AmountCents, &p.Currency, &p.Status, &p.Reason, &p.IdempotencyKey, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Payment{}, shared.ErrNotFound
		}
		return model.Payment{}, err
	}
	return p, nil
}