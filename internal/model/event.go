package model

import "time"

type Event struct {
	ID             string    `json:"id"`
	SubscriptionID string    `json:"subscription_id"`
	Type           string    `json:"type"`
	Payload        []byte    `json:"payload"`
	CreatedAt      time.Time `json:"created_at"`
}
