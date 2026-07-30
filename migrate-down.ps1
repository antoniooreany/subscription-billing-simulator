$ErrorActionPreference = 'Stop'
$env:DATABASE_URL = if ($env:DATABASE_URL) { $env:DATABASE_URL } else { 'postgres://postgres:postgres@localhost:5432/billing?sslmode=disable' }
psql $env:DATABASE_URL -f migrations/001_init.down.sql