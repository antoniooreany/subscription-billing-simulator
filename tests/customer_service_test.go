package tests

import (
	"context"
	"testing"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/antoniooreany/subscription-billing-simulator/internal/service"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
)

type customerRepoStub struct {
	createFn func(ctx context.Context, c model.Customer) (model.Customer, error)
}

func (r *customerRepoStub) Create(ctx context.Context, c model.Customer) (model.Customer, error) {
	return r.createFn(ctx, c)
}

func (r *customerRepoStub) GetByID(ctx context.Context, id string) (model.Customer, error) {
	return model.Customer{}, nil
}

func TestCustomerCreateRejectsEmptyFields(t *testing.T) {
	svc := service.NewCustomerService(nil)

	_, err := svc.Create(context.Background(), "", "Anton")
	if err != shared.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput for empty email, got %v", err)
	}

	_, err = svc.Create(context.Background(), "anton@example.com", "")
	if err != shared.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput for empty name, got %v", err)
	}
}
