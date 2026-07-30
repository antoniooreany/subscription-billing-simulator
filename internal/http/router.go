package http

import (
	stdhttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handler) stdhttp.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", h.Health)
	r.Post("/customers", h.CreateCustomer)
	r.Post("/subscriptions", h.CreateSubscription)
	r.Get("/subscriptions/{id}", h.GetSubscription)
	r.Get("/subscriptions/{id}/events", h.ListEvents)
	r.Post("/payments/fail", h.FailPayment)
	r.Post("/retries/run", h.RunRetry)
	return r
}