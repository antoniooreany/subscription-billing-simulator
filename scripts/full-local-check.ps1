$ErrorActionPreference = "Stop"

Write-Host "== Full local check =="

Write-Host "Cleaning up existing process on port 8080..."
$existing = Get-NetTCPConnection -LocalPort 8080 -State Listen -ErrorAction SilentlyContinue |
  Select-Object -ExpandProperty OwningProcess -Unique

if ($existing) {
  foreach ($existingPid in $existing) {
    if ($existingPid -and $existingPid -ne 0) {
      Write-Host "Stopping process on 8080: PID $existingPid"
      Stop-Process -Id $existingPid -Force -ErrorAction SilentlyContinue
    }
  }

  Start-Sleep -Seconds 2
}

docker compose up -d db

Write-Host "Waiting for PostgreSQL readiness..."
$pgReady = $false
$pgAttempts = 40

for ($i = 0; $i -lt $pgAttempts; $i++) {
  try {
    docker exec billing_db pg_isready -U postgres -d billing *> $null
    if ($LASTEXITCODE -eq 0) {
      $pgReady = $true
      break
    }
  } catch {
  }

  Start-Sleep -Seconds 2
}

if (-not $pgReady) {
  Write-Host ""
  Write-Host "PostgreSQL did not become ready in time."
  docker compose logs db --tail 100
  throw "PostgreSQL readiness check failed"
}

Write-Host "PostgreSQL is ready."

$env:APP_PORT = "8080"
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/billing?sslmode=disable"
$env:AUTO_MIGRATE = "true"

$projectRoot = Split-Path -Parent $PSScriptRoot
$stdoutLog = Join-Path $projectRoot "api.stdout.log"
$stderrLog = Join-Path $projectRoot "api.stderr.log"

if (Test-Path $stdoutLog) { Remove-Item $stdoutLog -Force }
if (Test-Path $stderrLog) { Remove-Item $stderrLog -Force }

$apiProcess = $null

try {
  Write-Host "Starting API..."
  $apiProcess = Start-Process powershell `
    -ArgumentList "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", ".\scripts\run.ps1" `
    -WorkingDirectory $projectRoot `
    -RedirectStandardOutput $stdoutLog `
    -RedirectStandardError $stderrLog `
    -PassThru

  Write-Host "Waiting for API health..."
  $healthy = $false
  $maxAttempts = 30

  for ($attempt = 0; $attempt -lt $maxAttempts; $attempt++) {
    try {
      $health = Invoke-RestMethod "http://localhost:8080/health"
      if ($health.status -eq "ok") {
        $healthy = $true
        break
      }
    } catch {
    }

    Start-Sleep -Seconds 2
  }

  if (-not $healthy) {
    Write-Host ""
    Write-Host "API stdout:"
    if (Test-Path $stdoutLog) { Get-Content $stdoutLog }

    Write-Host ""
    Write-Host "API stderr:"
    if (Test-Path $stderrLog) { Get-Content $stderrLog }

    Write-Host ""
    Write-Host "DB logs:"
    docker compose logs db --tail 100

    throw "API did not become healthy in time"
  }

  Write-Host "API is healthy."
  Write-Host "Waiting for database schema readiness..."

  $schemaReady = $false
  $schemaAttempts = 40

  for ($i = 0; $i -lt $schemaAttempts; $i++) {
    try {
      $result = docker exec billing_db psql -U postgres -d billing -t -c "select to_regclass('public.customers');" 2>$null
      if ($result -match "customers") {
        $schemaReady = $true
        break
      }
    } catch {
    }

    Start-Sleep -Seconds 2
  }

  if (-not $schemaReady) {
    Write-Host ""
    Write-Host "API stdout:"
    if (Test-Path $stdoutLog) { Get-Content $stdoutLog }

    Write-Host ""
    Write-Host "API stderr:"
    if (Test-Path $stderrLog) { Get-Content $stderrLog }

    Write-Host ""
    Write-Host "DB logs:"
    docker compose logs db --tail 100

    throw "Database schema did not become ready in time"
  }

  Write-Host "Database schema is ready."
  Write-Host "Running smoke test..."
  & "$PSScriptRoot\smoke-test.ps1"

  Write-Host ""
  Write-Host "Full local check passed."

  Write-Host ""
  Write-Host "API stdout:"
  if (Test-Path $stdoutLog) { Get-Content $stdoutLog }

  Write-Host ""
  Write-Host "API stderr:"
  if (Test-Path $stderrLog) { Get-Content $stderrLog }
}
finally {
  if ($apiProcess -and -not $apiProcess.HasExited) {
    Write-Host ""
    Write-Host "Stopping API process: $($apiProcess.Id)"
    Stop-Process -Id $apiProcess.Id -Force -ErrorAction SilentlyContinue
  }
}
