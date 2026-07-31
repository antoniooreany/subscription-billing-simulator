$ErrorActionPreference = 'Stop'
$env:APP_PORT = if ($env:APP_PORT) { $env:APP_PORT } else { '8080' }
$env:DATABASE_URL = if ($env:DATABASE_URL) { $env:DATABASE_URL } else { 'postgres://postgres:postgres@localhost:5432/billing?sslmode=disable' }
$env:AUTO_MIGRATE = if ($env:AUTO_MIGRATE) { $env:AUTO_MIGRATE } else { 'true' }
go run ./cmd/api