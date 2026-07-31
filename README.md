# Subscription Billing Simulator

> For reviewers
>
> This repository contains a small Go REST API that simulates subscription billing flows with PostgreSQL and Docker Compose. To review the project quickly, follow the “Quick start” and “Reviewer checklist” sections below.

## Release status

Current release: `v0.1.0`

### Included in v0.1.0
- Go REST API for subscription billing workflow simulation.
- PostgreSQL-backed data model for customers, subscriptions, payments, retry attempts, and events.
- Auto-applied SQL migrations.
- Docker Compose local environment.
- Idempotent payment failure and retry endpoints.
- End-to-end smoke test and review-ready local verification flow.

## Quick start

### Prerequisites
- Go
- Docker and Docker Compose
- PowerShell or a POSIX shell

### Run locally
```bash
docker compose up -d db
```

Set environment variables:

```bash
APP_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/billing?sslmode=disable
AUTO_MIGRATE=true
```

Start the API:

```bash
go run ./cmd/api
```

### Verify
Health check:

```bash
curl http://localhost:8080/health
```

Run smoke test:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\smoke-test.ps1
```

A small Go REST API that simulates a subscription billing workflow: customer creation, subscription lifecycle, failed payments, retry attempts, and event history.

## Highlights

- Implemented a Go-based REST API simulating subscription billing and recovery flows.
- Modeled customers, subscriptions, payments, retry attempts, and event history in PostgreSQL.
- Added idempotent payment failure and retry endpoints using idempotency keys.
- Built a reproducible local developer workflow with Docker Compose, startup migrations, and PowerShell verification scripts.
- Organized the codebase into clear handler, service, repository, and model layers.

## Features

- Create customers.
- Create subscriptions.
- Fetch subscription details.
- Register failed payments.
- Run retry attempts.
- Reactivate subscriptions after retries.
- View subscription event history.
- Apply SQL migrations automatically on startup.
- Run locally with PostgreSQL via Docker Compose.

## Stack

- Go
- Chi router
- PostgreSQL
- pgx
- Docker Compose
- PowerShell helpers for Windows

## Architecture

The API is organized into clear layers:

- Handlers receive HTTP requests, decode input, and return JSON responses.
- Services contain billing workflow logic such as failed payments, retry attempts, and subscription reactivation.
- Repositories encapsulate PostgreSQL access through pgx.
- Models define the domain entities shared across the application.

## Design decisions

- PostgreSQL is the source of truth for customers, subscriptions, payments, retry attempts, and events.
- Idempotency keys are required for `payments/fail` and `retries/run` to avoid duplicate writes during retries.
- Failed payments move subscriptions to `past_due`.
- The second retry attempt reactivates the subscription and appends a `subscription_reactivated` event.
- Event history is stored in the database so the billing flow can be inspected after each run.

## Project structure

```text
subscription-billing-simulator/
├── cmd/
├── internal/
├── migrations/
├── scripts/
│   ├── run.ps1
│   ├── migrate-up.ps1
│   ├── migrate-down.ps1
│   ├── smoke-test.ps1
│   └── full-local-check.ps1
├── tests/
├── docker-compose.yml
├── go.mod
└── README.md
```

## Run locally

### 1. Start PostgreSQL

```powershell
docker compose up -d db
```

### 2. Set environment variables

```powershell
$env:APP_PORT="8080"
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/billing?sslmode=disable"
$env:AUTO_MIGRATE="true"
```

### 3. Start the API

```powershell
.\scripts\run.ps1
```

If startup succeeds, the server should print:

```text
server listening on :8080
```

## Health check

```powershell
Invoke-RestMethod http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

## How to verify

Use the helper scripts instead of running the whole flow manually.

### Smoke test

```powershell
.\scripts\smoke-test.ps1
```

### Full local check

This command starts PostgreSQL, starts the API, waits for `/health`, and then runs the smoke test.

```powershell
.\scripts\full-local-check.ps1
```

### What the smoke test covers

- Health check.
- Customer creation.
- Subscription creation.
- Failed payment registration.
- Retry attempt 1.
- Retry attempt 2.
- Subscription reactivation.
- Event history fetch.

## Expected flow

The expected lifecycle in the verification flow is:

1. A customer is created.
2. A subscription is created with status `active`.
3. A failed payment is registered.
4. A first retry attempt is processed.
5. A second retry attempt is processed.
6. The subscription is reactivated and ends in status `active`.
7. Event history includes:
   - `subscription_created`
   - `payment_failed`
   - `retry_processed`
   - `retry_processed`
   - `subscription_reactivated`

## Reviewer checklist

You can use this short checklist to verify the project locally.

1. Clone and checkout the release:
   ```bash
   git clone https://github.com/antoniooreany/subscription-billing-simulator.git
   cd subscription-billing-simulator
   git checkout v0.1.0
   ```

2. Start PostgreSQL:
   ```bash
   docker compose up -d db
   ```

3. Set environment variables:
   ```bash
   APP_PORT=8080
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/billing?sslmode=disable
   AUTO_MIGRATE=true
   ```

4. Start the API:
   ```bash
   go run ./cmd/api
   ```

5. Run health check:
   ```bash
   curl http://localhost:8080/health
   ```
   Expected result: `{"status":"ok"}`

6. Run the smoke test:
   ```powershell
   .\scripts\smoke-test.ps1
   ```
   Expected result: the script completes without errors and verifies the main billing flow.

7. Optional full local verification:
   ```powershell
   .\scripts\full-local-check.ps1
   ```
   Expected result: PostgreSQL and the API are started, `/health` passes, and the smoke test runs end-to-end.

## Example API requests

### Create customer

```powershell
$uniqueEmail = "anton+$(Get-Date -Format 'yyyyMMddHHmmss')@example.com"

$customerBody = @{
  email = $uniqueEmail
  name = "Anton Gorshkov"
} | ConvertTo-Json

Invoke-RestMethod `
  -Uri "http://localhost:8080/customers" `
  -Method POST `
  -ContentType "application/json" `
  -Body $customerBody
```

### Create subscription

```powershell
$subscriptionBody = @{
  customer_id = $customer.id
  plan_code = "basic-monthly"
  amount_cents = 990
  currency = "EUR"
} | ConvertTo-Json

Invoke-RestMethod `
  -Uri "http://localhost:8080/subscriptions" `
  -Method POST `
  -ContentType "application/json" `
  -Body $subscriptionBody
```

### Fail payment

```powershell
$paymentFailBody = @{
  subscription_id = $subscription.id
  amount_cents = 990
  currency = "EUR"
  reason = "card_declined"
  idempotency_key = "fail-$($subscription.id)-001"
} | ConvertTo-Json

Invoke-RestMethod `
  -Uri "http://localhost:8080/payments/fail" `
  -Method POST `
  -ContentType "application/json" `
  -Body $paymentFailBody
```

### Run retry

```powershell
$retryBody = @{
  subscription_id = $subscription.id
  idempotency_key = "retry-$($subscription.id)-001"
} | ConvertTo-Json

Invoke-RestMethod `
  -Uri "http://localhost:8080/retries/run" `
  -Method POST `
  -ContentType "application/json" `
  -Body $retryBody
```

### Get subscription

```powershell
Invoke-RestMethod "http://localhost:8080/subscriptions/$($subscription.id)"
```

### Get events

```powershell
Invoke-RestMethod "http://localhost:8080/subscriptions/$($subscription.id)/events"
```

## Reset local database

```powershell
docker compose down -v
.\scripts\full-local-check.ps1
```

## Notes

- SQL migrations are applied automatically on startup.
- The initial SQL scaffold is suitable for local development and repeated verification runs.
- Use a unique email when creating a customer to avoid duplicate key errors on `customers.email`.
- `payments/fail` and `retries/run` require `idempotency_key` in the request body.
- The Docker Compose warning about `version` being obsolete can be removed by deleting the `version` field from `docker-compose.yml`.

## Current status

Working:

- health check
- customer creation
- subscription creation
- subscription fetch
- failed payment flow
- retry flow
- subscription reactivation
- event history fetch

Next improvements:

- improve validation error messages
- return more specific API errors for missing fields
- add more integration tests
- refine retry policy and billing rules
- add seed data or example fixtures
## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
