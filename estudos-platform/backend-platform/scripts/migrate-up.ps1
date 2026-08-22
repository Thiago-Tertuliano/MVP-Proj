# Aplica as migrations no Postgres (golang-migrate).
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root
go run ./cmd/migrate
