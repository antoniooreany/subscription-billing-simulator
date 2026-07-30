$ErrorActionPreference = "Stop"

function Step($name) {
  Write-Host ""
  Write-Host "== $name =="
}

function Fail-Step($name, $details) {
  Write-Host ""
  Write-Host "FAILED: $name" -ForegroundColor Red
  if ($details) {
    Write-Host $details -ForegroundColor DarkRed
  }
  throw "Smoke test failed at step: $name"
}

function Invoke-JsonPost($name, $uri, $bodyObject) {
  try {
    return Invoke-RestMethod `
      -Uri $uri `
      -Method POST `
      -ContentType "application/json" `
      -Body ($bodyObject | ConvertTo-Json -Depth 10)
  }
  catch {
    $message = $_.Exception.Message
    if ($_.ErrorDetails -and $_.ErrorDetails.Message) {
      $message = $_.ErrorDetails.Message
    }
    Fail-Step $name $message
  }
}

Write-Host "== Subscription Billing Simulator smoke test =="

Step "Health"
try {
  $health = Invoke-RestMethod "http://localhost:8080/health"
  if ($health.status -ne "ok") {
    Fail-Step "Health" "Expected status 'ok' but got '$($health.status)'"
  }
  Write-Host "Health: $($health.status)"
}
catch {
  $message = $_.Exception.Message
  if ($_.ErrorDetails -and $_.ErrorDetails.Message) {
    $message = $_.ErrorDetails.Message
  }
  Fail-Step "Health" $message
}

Step "Create customer"
$uniqueEmail = "anton+$(Get-Date -Format 'yyyyMMddHHmmssfff')@example.com"
$customer = Invoke-JsonPost `
  "Create customer" `
  "http://localhost:8080/customers" `
  @{
    email = $uniqueEmail
    name  = "Anton Gorshkov"
  }

if (-not $customer.id) {
  Fail-Step "Create customer" "Response did not contain customer id"
}
Write-Host "Customer ID: $($customer.id)"

Step "Create subscription"
$subscription = Invoke-JsonPost `
  "Create subscription" `
  "http://localhost:8080/subscriptions" `
  @{
    customer_id  = $customer.id
    plan_code    = "basic-monthly"
    amount_cents = 990
    currency     = "EUR"
  }

if (-not $subscription.id) {
  Fail-Step "Create subscription" "Response did not contain subscription id"
}
if ($subscription.status -ne "active") {
  Fail-Step "Create subscription" "Expected status 'active' but got '$($subscription.status)'"
}
Write-Host "Subscription ID: $($subscription.id)"
Write-Host "Initial status: $($subscription.status)"

Step "Fail payment"
$paymentFailResult = Invoke-JsonPost `
  "Fail payment" `
  "http://localhost:8080/payments/fail" `
  @{
    subscription_id = $subscription.id
    amount_cents    = 990
    currency        = "EUR"
    reason          = "card_declined"
    idempotency_key = "fail-$($subscription.id)-001"
  }

if ($paymentFailResult.status -ne "failed") {
  Fail-Step "Fail payment" "Expected payment status 'failed' but got '$($paymentFailResult.status)'"
}
Write-Host "Failed payment status: $($paymentFailResult.status)"

Step "Retry attempt 1"
$retryResult1 = Invoke-JsonPost `
  "Retry attempt 1" `
  "http://localhost:8080/retries/run" `
  @{
    subscription_id = $subscription.id
    idempotency_key = "retry-$($subscription.id)-001"
  }

if ($retryResult1.attempt_number -ne 1) {
  Fail-Step "Retry attempt 1" "Expected attempt_number 1 but got '$($retryResult1.attempt_number)'"
}
Write-Host "Retry 1 attempt number: $($retryResult1.attempt_number)"

Step "Retry attempt 2"
$retryResult2 = Invoke-JsonPost `
  "Retry attempt 2" `
  "http://localhost:8080/retries/run" `
  @{
    subscription_id = $subscription.id
    idempotency_key = "retry-$($subscription.id)-002"
  }

if ($retryResult2.attempt_number -ne 2) {
  Fail-Step "Retry attempt 2" "Expected attempt_number 2 but got '$($retryResult2.attempt_number)'"
}
Write-Host "Retry 2 attempt number: $($retryResult2.attempt_number)"

Step "Fetch final subscription"
try {
  $subscriptionDetails = Invoke-RestMethod "http://localhost:8080/subscriptions/$($subscription.id)"
}
catch {
  $message = $_.Exception.Message
  if ($_.ErrorDetails -and $_.ErrorDetails.Message) {
    $message = $_.ErrorDetails.Message
  }
  Fail-Step "Fetch final subscription" $message
}

if ($subscriptionDetails.status -ne "active") {
  Fail-Step "Fetch final subscription" "Expected final status 'active' but got '$($subscriptionDetails.status)'"
}
Write-Host "Final subscription status: $($subscriptionDetails.status)"

Step "Fetch events"
try {
  $events = Invoke-RestMethod "http://localhost:8080/subscriptions/$($subscription.id)/events"
}
catch {
  $message = $_.Exception.Message
  if ($_.ErrorDetails -and $_.ErrorDetails.Message) {
    $message = $_.ErrorDetails.Message
  }
  Fail-Step "Fetch events" $message
}

$eventTypes = @($events | ForEach-Object { $_.type })
$expectedTypes = @(
  "subscription_created",
  "payment_failed",
  "retry_processed",
  "retry_processed",
  "subscription_reactivated"
)

if ($eventTypes.Count -ne 5) {
  Fail-Step "Fetch events" "Expected 5 events but got $($eventTypes.Count)"
}

for ($i = 0; $i -lt $expectedTypes.Count; $i++) {
  if ($eventTypes[$i] -ne $expectedTypes[$i]) {
    Fail-Step "Fetch events" "Expected event[$i] '$($expectedTypes[$i])' but got '$($eventTypes[$i])'"
  }
}

Write-Host "Events count: $($eventTypes.Count)"

Write-Host ""
Write-Host "===== SUMMARY ====="
Write-Host "Health: $($health.status)"
Write-Host "Customer ID: $($customer.id)"
Write-Host "Subscription ID: $($subscription.id)"
Write-Host "Initial status: $($subscription.status)"
Write-Host "Failed payment status: $($paymentFailResult.status)"
Write-Host "Retry 1 attempt number: $($retryResult1.attempt_number)"
Write-Host "Retry 2 attempt number: $($retryResult2.attempt_number)"
Write-Host "Final subscription status: $($subscriptionDetails.status)"
Write-Host "Events count: $($eventTypes.Count)"
Write-Host ""
Write-Host "Smoke test passed successfully." -ForegroundColor Green
