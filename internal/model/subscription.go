package model

import "time"

type SubscriptionStatus string

const (
	SubscriptionActive   SubscriptionStatus = "active"
	SubscriptionPastDue  SubscriptionStatus = "past_due"
	SubscriptionCanceled SubscriptionStatus = "canceled"
)

type Subscription struct {
	ID          string             `json:"id"`
	CustomerID  string             `json:"customer_id"`
	PlanCode    string             `json:"plan_code"`
	AmountCents int                `json:"amount_cents"`
	Currency    string             `json:"currency"`
	Status      SubscriptionStatus `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
}
