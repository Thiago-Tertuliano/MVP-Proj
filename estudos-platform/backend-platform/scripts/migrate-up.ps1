# Aplica as migrations no container estudos-postgres (PowerShell / Windows).
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$mig = Join-Path $root "internal\infrastructure\persistence\postgres\migration"

$files = @(
    "0001_create_usuario.up.sql",
    "0002_create_artigos.up.sql",
    "0003_create_trilhas.up.sql"
)

foreach ($f in $files) {
    $path = Join-Path $mig $f
    Write-Host ">> $f"
    Get-Content $path -Raw | docker exec -i estudos-postgres psql -U estudos -d estudos_platform
}

Write-Host ">> tabelas:"
docker exec estudos-postgres psql -U estudos -d estudos_platform -c "\dt"
