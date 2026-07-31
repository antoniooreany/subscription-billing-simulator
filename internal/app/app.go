package app

import (
	"context"
	stdhttp "net/http"
	"path/filepath"

	"github.com/antoniooreany/subscription-billing-simulator/internal/config"
	"github.com/antoniooreany/subscription-billing-simulator/internal/db"
	httpapi "github.com/antoniooreany/subscription-billing-simulator/internal/http"
	"github.com/antoniooreany/subscription-billing-simulator/internal/repository"
	"github.com/antoniooreany/subscription-billing-simulator/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config config.Config
	DB     *pgxpool.Pool
	Router stdhttp.Handler
}

func New() (*App, error) {
	cfg := config.Load()
	pool, err := db.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if cfg.AutoMigrate {
		if err := db.ApplyMigrations(context.Background(), pool, filepath.Join("migrations")); err != nil {
			pool.Close()
			return nil, err
		}
	}

	customerRepo := repository.NewCustomerRepository(pool)
	subscriptionRepo := repository.NewSubscriptionRepository(pool)
	paymentRepo := repository.NewPaymentRepository(pool)
	retryRepo := repository.NewRetryRepository(pool)
	eventRepo := repository.NewEventRepository(pool)

	customerService := service.NewCustomerService(customerRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, eventRepo)
	paymentService := service.NewPaymentService(paymentRepo, subscriptionRepo, eventRepo)
	retryService := service.NewRetryService(retryRepo, subscriptionRepo, eventRepo)

	handler := httpapi.NewHandler(customerService, subscriptionService, paymentService, retryService)
	router := httpapi.NewRouter(handler)

	return &App{Config: cfg, DB: pool, Router: router}, nil
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
