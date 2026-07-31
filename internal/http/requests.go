package http

type CreateCustomerRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CreateSubscriptionRequest struct {
	CustomerID  string `json:"customer_id"`
	PlanCode    string `json:"plan_code"`
	AmountCents int    `json:"amount_cents"`
	Currency    string `json:"currency"`
}

type FailPaymentRequest struct {
	SubscriptionID string `json:"subscription_id"`
	AmountCents    int    `json:"amount_cents"`
	Currency       string `json:"currency"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

type RunRetryRequest struct {
	SubscriptionID string `json:"subscription_id"`
	IdempotencyKey string `json:"idempotency_key"`
}
