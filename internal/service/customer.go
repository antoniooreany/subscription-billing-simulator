package service

import (
	"context"
	"strings"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/antoniooreany/subscription-billing-simulator/internal/repository"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(ctx context.Context, email, name string) (model.Customer, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(name) == "" {
		return model.Customer{}, shared.ErrInvalidInput
	}
	c := model.Customer{ID: shared.NewID("cust"), Email: email, Name: name}
	return s.repo.Create(ctx, c)
}
