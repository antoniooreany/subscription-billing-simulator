package model

import "time"

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentFailed  PaymentStatus = "failed"
	PaymentPaid    PaymentStatus = "paid"
)

type Payment struct {
	ID             string        `json:"id"`
	SubscriptionID string        `json:"subscription_id"`
	AmountCents    int           `json:"amount_cents"`
	Currency       string        `json:"currency"`
	Status         PaymentStatus `json:"status"`
	Reason         string        `json:"reason,omitempty"`
	IdempotencyKey string        `json:"idempotency_key,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
}
