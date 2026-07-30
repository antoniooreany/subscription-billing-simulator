package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/antoniooreany/subscription-billing-simulator/internal/repository"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
)

type RetryService struct {
	retries       *repository.RetryRepository
	subscriptions *repository.SubscriptionRepository
	events        *repository.EventRepository
}

func NewRetryService(retries *repository.RetryRepository, subscriptions *repository.SubscriptionRepository, events *repository.EventRepository) *RetryService {
	return &RetryService{retries: retries, subscriptions: subscriptions, events: events}
}

func (s *RetryService) Run(ctx context.Context, subscriptionID, key string) (model.RetryAttempt, error) {
	if strings.TrimSpace(subscriptionID) == "" || strings.TrimSpace(key) == "" {
		return model.RetryAttempt{}, shared.ErrInvalidInput
	}
	if existing, err := s.retries.GetByIdempotencyKey(ctx, key); err == nil {
		return existing, nil
	}
	if _, err := s.subscriptions.GetByID(ctx, subscriptionID); err != nil {
		return model.RetryAttempt{}, err
	}
	count, err := s.retries.CountBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return model.RetryAttempt{}, err
	}
	attempt := model.RetryAttempt{ID: shared.NewID("retry"), SubscriptionID: subscriptionID, AttemptNumber: count + 1, Status: model.RetryProcessed, IdempotencyKey: key}
	created, err := s.retries.Create(ctx, attempt)
	if err != nil {
		return model.RetryAttempt{}, err
	}
	payload := fmt.Sprintf(`{"attempt_number":%d,"status":"processed"}`, created.AttemptNumber)
	_, _ = s.events.Create(ctx, model.Event{ID: shared.NewID("evt"), SubscriptionID: subscriptionID, Type: "retry_processed", Payload: []byte(payload)})
	if created.AttemptNumber >= 2 {
		_ = s.subscriptions.UpdateStatus(ctx, subscriptionID, model.SubscriptionActive)
		_, _ = s.events.Create(ctx, model.Event{ID: shared.NewID("evt"), SubscriptionID: subscriptionID, Type: "subscription_reactivated", Payload: []byte(`{"status":"active"}`)})
	}
	return created, nil
}
