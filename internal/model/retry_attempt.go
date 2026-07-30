package model

import "time"

type RetryStatus string

const (
	RetryScheduled RetryStatus = "scheduled"
	RetryProcessed RetryStatus = "processed"
	RetryResolved  RetryStatus = "resolved"
)

type RetryAttempt struct {
	ID             string      `json:"id"`
	SubscriptionID string      `json:"subscription_id"`
	AttemptNumber  int         `json:"attempt_number"`
	Status         RetryStatus `json:"status"`
	IdempotencyKey string      `json:"idempotency_key,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
}