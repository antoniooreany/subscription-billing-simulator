package tests

import (
	"context"
	"testing"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/antoniooreany/subscription-billing-simulator/internal/service"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
)

type paymentRepoFake struct {
	getByKeyFn func(ctx context.Context, key string) (model.Payment, error)
	createFn   func(ctx context.Context, p model.Payment) (model.Payment, error)
}

func (f *paymentRepoFake) GetByIdempotencyKey(ctx context.Context, key string) (model.Payment, error) {
	return f.getByKeyFn(ctx, key)
}

func (f *paymentRepoFake) Create(ctx context.Context, p model.Payment) (model.Payment, error) {
	return f.createFn(ctx, p)
}

type subscriptionRepoFake struct {
	getByIDFn      func(ctx context.Context, id string) (model.Subscription, error)
	updateStatusFn func(ctx context.Context, id string, status model.SubscriptionStatus) error
	createFn       func(ctx context.Context, s model.Subscription) (model.Subscription, error)
}

func (f *subscriptionRepoFake) Create(ctx context.Context, s model.Subscription) (model.Subscription, error) {
	if f.createFn != nil {
		return f.createFn(ctx, s)
	}
	return s, nil
}

func (f *subscriptionRepoFake) GetByID(ctx context.Context, id string) (model.Subscription, error) {
	return f.getByIDFn(ctx, id)
}

func (f *subscriptionRepoFake) UpdateStatus(ctx context.Context, id string, status model.SubscriptionStatus) error {
	return f.updateStatusFn(ctx, id, status)
}

type eventRepoFake struct {
	createCalls int
	createFn    func(ctx context.Context, e model.Event) (model.Event, error)
	listFn      func(ctx context.Context, subscriptionID string) ([]model.Event, error)
}

func (f *eventRepoFake) Create(ctx context.Context, e model.Event) (model.Event, error) {
	f.createCalls++
	if f.createFn != nil {
		return f.createFn(ctx, e)
	}
	return e, nil
}

func (f *eventRepoFake) ListBySubscriptionID(ctx context.Context, subscriptionID string) ([]model.Event, error) {
	if f.listFn != nil {
		return f.listFn(ctx, subscriptionID)
	}
	return nil, nil
}

func TestFailPaymentCreatesFailedPaymentAndMarksSubscriptionPastDue(t *testing.T) {
	var updatedStatus model.SubscriptionStatus

	svc := service.NewPaymentService(
		&paymentRepoFake{
			getByKeyFn: func(ctx context.Context, key string) (model.Payment, error) {
				return model.Payment{}, shared.ErrNotFound
			},
			createFn: func(ctx context.Context, p model.Payment) (model.Payment, error) {
				return p, nil
			},
		},
		&subscriptionRepoFake{
			getByIDFn: func(ctx context.Context, id string) (model.Subscription, error) {
				return model.Subscription{ID: id, Status: model.SubscriptionActive}, nil
			},
			updateStatusFn: func(ctx context.Context, id string, status model.SubscriptionStatus) error {
				updatedStatus = status
				return nil
			},
		},
		&eventRepoFake{},
	)

	res, err := svc.FailPayment(context.Background(), "sub_123", 990, "EUR", "card_declined", "fail-sub_123-001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Status != model.PaymentFailed {
		t.Fatalf("expected payment status failed, got %s", res.Status)
	}

	if updatedStatus != model.SubscriptionPastDue {
		t.Fatalf("expected subscription status to be updated to past_due, got %s", updatedStatus)
	}
}
