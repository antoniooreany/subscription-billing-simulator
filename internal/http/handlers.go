package http

import (
	"encoding/json"
	stdhttp "net/http"

	"github.com/antoniooreany/subscription-billing-simulator/internal/service"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	customers     *service.CustomerService
	subscriptions *service.SubscriptionService
	payments      *service.PaymentService
	retries       *service.RetryService
}

func NewHandler(customers *service.CustomerService, subscriptions *service.SubscriptionService, payments *service.PaymentService, retries *service.RetryService) *Handler {
	return &Handler{customers: customers, subscriptions: subscriptions, payments: payments, retries: retries}
}

func (h *Handler) Health(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
	writeJSON(w, stdhttp.StatusOK, HealthResponse{Status: "ok"})
}

func (h *Handler) CreateCustomer(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid json")
		return
	}
	res, err := h.customers.Create(r.Context(), req.Email, req.Name)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusCreated, res)
}

func (h *Handler) CreateSubscription(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid json")
		return
	}
	res, err := h.subscriptions.Create(r.Context(), req.CustomerID, req.PlanCode, req.AmountCents, req.Currency)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusCreated, res)
}

func (h *Handler) GetSubscription(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	res, err := h.subscriptions.GetByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, res)
}

func (h *Handler) FailPayment(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req FailPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid json")
		return
	}
	res, err := h.payments.FailPayment(r.Context(), req.SubscriptionID, req.AmountCents, req.Currency, req.Reason, req.IdempotencyKey)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, res)
}

func (h *Handler) RunRetry(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req RunRetryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid json")
		return
	}
	res, err := h.retries.Run(r.Context(), req.SubscriptionID, req.IdempotencyKey)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, res)
}

func (h *Handler) ListEvents(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	res, err := h.subscriptions.Events(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, res)
}

func writeJSON(w stdhttp.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w stdhttp.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func handleError(w stdhttp.ResponseWriter, err error) {
	switch err {
	case shared.ErrInvalidInput:
		writeError(w, stdhttp.StatusBadRequest, err.Error())
	case shared.ErrNotFound:
		writeError(w, stdhttp.StatusNotFound, err.Error())
	case shared.ErrConflict:
		writeError(w, stdhttp.StatusConflict, err.Error())
	default:
		writeError(w, stdhttp.StatusInternalServerError, "internal server error")
	}
}
