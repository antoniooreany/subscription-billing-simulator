package service

import (
	"context"
	"strings"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s model.Subscription) (model.Subscription, error)
	GetByID(ctx context.Context, id string) (model.Subscription, error)
	UpdateStatus(ctx context.Context, id string, status model.SubscriptionStatus) error
}

type EventRepository interface {
	Create(ctx context.Context, e model.Event) (model.Event, error)
	ListBySubscriptionID(ctx context.Context, subscriptionID string) ([]model.Event, error)
}

type SubscriptionService struct {
	repo      SubscriptionRepository
	eventRepo EventRepository
}

func NewSubscriptionService(repo SubscriptionRepository, eventRepo EventRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo, eventRepo: eventRepo}
}

func (s *SubscriptionService) Create(ctx context.Context, customerID, planCode string, amountCents int, currency string) (model.Subscription, error) {
	if strings.TrimSpace(customerID) == "" || strings.TrimSpace(planCode) == "" || strings.TrimSpace(currency) == "" || amountCents <= 0 {
		return model.Subscription{}, shared.ErrInvalidInput
	}

	sub := model.Subscription{
		ID:          shared.NewID("sub"),
		CustomerID:  customerID,
		PlanCode:    planCode,
		AmountCents: amountCents,
		Currency:    currency,
		Status:      model.SubscriptionActive,
	}

	created, err := s.repo.Create(ctx, sub)
	if err != nil {
		return model.Subscription{}, err
	}

	_, _ = s.eventRepo.Create(ctx, model.Event{
		ID:             shared.NewID("evt"),
		SubscriptionID: created.ID,
		Type:           "subscription_created",
		Payload:        []byte(`{"status":"active"}`),
	})

	return created, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id string) (model.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubscriptionService) Events(ctx context.Context, subscriptionID string) ([]model.Event, error) {
	return s.eventRepo.ListBySubscriptionID(ctx, subscriptionID)
}
