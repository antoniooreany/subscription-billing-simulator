package repository

import (
	"context"
	"errors"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RetryRepository struct{ db *pgxpool.Pool }

func NewRetryRepository(db *pgxpool.Pool) *RetryRepository { return &RetryRepository{db: db} }

func (r *RetryRepository) Create(ctx context.Context, ra model.RetryAttempt) (model.RetryAttempt, error) {
	_, err := r.db.Exec(ctx, `INSERT INTO retry_attempts (id, subscription_id, attempt_number, status, idempotency_key) VALUES ($1, $2, $3, $4, $5)`, ra.ID, ra.SubscriptionID, ra.AttemptNumber, ra.Status, ra.IdempotencyKey)
	return ra, err
}

func (r *RetryRepository) CountBySubscriptionID(ctx context.Context, subscriptionID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM retry_attempts WHERE subscription_id = $1`, subscriptionID).Scan(&count)
	return count, err
}

func (r *RetryRepository) GetByIdempotencyKey(ctx context.Context, key string) (model.RetryAttempt, error) {
	var ra model.RetryAttempt
	err := r.db.QueryRow(ctx, `SELECT id, subscription_id, attempt_number, status, idempotency_key, created_at FROM retry_attempts WHERE idempotency_key = $1`, key).Scan(&ra.ID, &ra.SubscriptionID, &ra.AttemptNumber, &ra.Status, &ra.IdempotencyKey, &ra.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RetryAttempt{}, shared.ErrNotFound
		}
		return model.RetryAttempt{}, err
	}
	return ra, nil
}
