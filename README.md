# Subscription Billing Simulator

A small Go REST API that simulates a subscription billing workflow: customers, subscriptions, failed payments, retry attempts, and event history.

## Features

- Create customers
- Create subscriptions
- Fetch subscription details
- View subscription event history
- Apply SQL migrations automatically on startup
- Run locally with PostgreSQL via Docker Compose

## Stack

- Go
- Chi router
- PostgreSQL
- pgx
- Docker Compose
- PowerShell helpers for Windows

## Project structure

```text
subscription-billing-simulator/
├── cmd/
├── internal/
├── migrations/
├── tests/
├── docker-compose.yml
├── run.ps1
├── migrate-up.ps1
├── migrate-down.ps1
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
.\run.ps1
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

Run the API and then execute the following PowerShell commands:

```powershell
# Health
$health = Invoke-RestMethod "http://localhost:8080/health"
$health

# Create customer
$uniqueEmail = "anton+$(Get-Date -Format 'yyyyMMddHHmmss')@example.com"

$customerBody = @{
  email = $uniqueEmail
  name  = "Anton Gorshkov"
} | ConvertTo-Json

$customer = Invoke-RestMethod `
  -Uri "http://localhost:8080/customers" `
  -Method POST `
  -ContentType "application/json" `
  -Body $customerBody

$customer

# Create subscription
$subscriptionBody = @{
  customer_id  = $customer.id
  plan_code    = "basic-monthly"
  amount_cents = 990
  currency     = "EUR"
} | ConvertTo-Json

$subscription = Invoke-RestMethod `
  -Uri "http://localhost:8080/subscriptions" `
  -Method POST `
  -ContentType "application/json" `
  -Body $subscriptionBody

$subscription

# Get subscription
$subscriptionDetails = Invoke-RestMethod "http://localhost:8080/subscriptions/$($subscription.id)"
$subscriptionDetails

# Get events
$events = Invoke-RestMethod "http://localhost:8080/subscriptions/$($subscription.id)/events"
$events

# Summary
Write-Host "Health: $($health.status)"
Write-Host "Customer ID: $($customer.id)"
Write-Host "Subscription ID: $($subscription.id)"
Write-Host "Subscription status: $($subscriptionDetails.status)"
```

## Example API flow

```powershell
$uniqueEmail = "anton+$(Get-Date -Format 'yyyyMMddHHmmss')@example.com"

$customerBody = @{
  email = $uniqueEmail
  name  = "Anton Gorshkov"
} | ConvertTo-Json

$customer = Invoke-RestMethod `
  -Uri "http://localhost:8080/customers" `
  -Method POST `
  -ContentType "application/json" `
  -Body $customerBody

$subscriptionBody = @{
  customer_id  = $customer.id
  plan_code    = "basic-monthly"
  amount_cents = 990
  currency     = "EUR"
} | ConvertTo-Json

$subscription = Invoke-RestMethod `
  -Uri "http://localhost:8080/subscriptions" `
  -Method POST `
  -ContentType "application/json" `
  -Body $subscriptionBody

Invoke-RestMethod "http://localhost:8080/subscriptions/$($subscription.id)"
Invoke-RestMethod "http://localhost:8080/subscriptions/$($subscription.id)/events"
```

## Reset local database

```powershell
docker compose down -v
docker compose up -d db
Start-Sleep -Seconds 15
.\run.ps1
```

## Current status

Working:
- health check
- customer creation
- subscription creation
- subscription fetch
- event history fetch

In progress:
- failed payment flow
- retry flow
- improved validation and error handling

## Notes

- SQL migrations are applied automatically on startup.
- The initial SQL scaffold is idempotent for local development.
- Use a unique email when creating a customer to avoid duplicate key errors on `customers.email`.
- The Docker Compose warning about `version` being obsolete can be removed by deleting the `version` field from `docker-compose.yml`.

## Next improvements

- Implement and verify failed payment flow
- Implement and verify retry flow
- Improve validation and API error messages
- Add unit tests for service logic
- Add integration tests for handlers and database
- Introduce proper migration tracking