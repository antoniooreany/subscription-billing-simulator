package tests

import (
	"context"
	"testing"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/antoniooreany/subscription-billing-simulator/internal/service"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
)

type retryRepoFake struct {
	getByKeyFn func(ctx context.Context, key string) (model.RetryAttempt, error)
	countFn    func(ctx context.Context, subscriptionID string) (int, error)
	createFn   func(ctx context.Context, ra model.RetryAttempt) (model.RetryAttempt, error)
}

func (f *retryRepoFake) GetByIdempotencyKey(ctx context.Context, key string) (model.RetryAttempt, error) {
	return f.getByKeyFn(ctx, key)
}

func (f *retryRepoFake) CountBySubscriptionID(ctx context.Context, subscriptionID string) (int, error) {
	return f.countFn(ctx, subscriptionID)
}

func (f *retryRepoFake) Create(ctx context.Context, ra model.RetryAttempt) (model.RetryAttempt, error) {
	return f.createFn(ctx, ra)
}

func TestRunRetrySecondAttemptReactivatesSubscription(t *testing.T) {
	var updatedStatus model.SubscriptionStatus
	events := &eventRepoFake{}

	svc := service.NewRetryService(
		&retryRepoFake{
			getByKeyFn: func(ctx context.Context, key string) (model.RetryAttempt, error) {
				return model.RetryAttempt{}, shared.ErrNotFound
			},
			countFn: func(ctx context.Context, subscriptionID string) (int, error) {
				return 1, nil
			},
			createFn: func(ctx context.Context, ra model.RetryAttempt) (model.RetryAttempt, error) {
				return ra, nil
			},
		},
		&subscriptionRepoFake{
			getByIDFn: func(ctx context.Context, id string) (model.Subscription, error) {
				return model.Subscription{ID: id, Status: model.SubscriptionPastDue}, nil
			},
			updateStatusFn: func(ctx context.Context, id string, status model.SubscriptionStatus) error {
				updatedStatus = status
				return nil
			},
		},
		events,
	)

	res, err := svc.Run(context.Background(), "sub_123", "retry-sub_123-002")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.AttemptNumber != 2 {
		t.Fatalf("expected attempt number 2, got %d", res.AttemptNumber)
	}

	if updatedStatus != model.SubscriptionActive {
		t.Fatalf("expected subscription to be reactivated, got %s", updatedStatus)
	}

	if events.createCalls < 2 {
		t.Fatalf("expected at least 2 event writes, got %d", events.createCalls)
	}
}
