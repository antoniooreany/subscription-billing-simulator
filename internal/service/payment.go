package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/antoniooreany/subscription-billing-simulator/internal/repository"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
)

type PaymentService struct {
	payments      *repository.PaymentRepository
	subscriptions *repository.SubscriptionRepository
	events        *repository.EventRepository
}

func NewPaymentService(payments *repository.PaymentRepository, subscriptions *repository.SubscriptionRepository, events *repository.EventRepository) *PaymentService {
	return &PaymentService{payments: payments, subscriptions: subscriptions, events: events}
}

func (s *PaymentService) FailPayment(ctx context.Context, subscriptionID string, amountCents int, currency, reason, key string) (model.Payment, error) {
	if strings.TrimSpace(subscriptionID) == "" || amountCents <= 0 || strings.TrimSpace(currency) == "" || strings.TrimSpace(key) == "" {
		return model.Payment{}, shared.ErrInvalidInput
	}
	if existing, err := s.payments.GetByIdempotencyKey(ctx, key); err == nil {
		return existing, nil
	}
	if _, err := s.subscriptions.GetByID(ctx, subscriptionID); err != nil {
		return model.Payment{}, err
	}
	payment := model.Payment{ID: shared.NewID("pay"), SubscriptionID: subscriptionID, AmountCents: amountCents, Currency: currency, Status: model.PaymentFailed, Reason: reason, IdempotencyKey: key}
	created, err := s.payments.Create(ctx, payment)
	if err != nil {
		return model.Payment{}, err
	}
	if err := s.subscriptions.UpdateStatus(ctx, subscriptionID, model.SubscriptionPastDue); err != nil {
		return model.Payment{}, err
	}
	payload := fmt.Sprintf(`{"reason":%q,"status":"failed"}`, reason)
	_, _ = s.events.Create(ctx, model.Event{ID: shared.NewID("evt"), SubscriptionID: subscriptionID, Type: "payment_failed", Payload: []byte(payload)})
	return created, nil
}