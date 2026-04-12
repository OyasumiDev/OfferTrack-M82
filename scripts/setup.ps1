# OfferTrack M82 - Setup de entorno (Windows PowerShell)
Write-Host "Verificando dependencias..." -ForegroundColor Cyan

if (!(Get-Command go    -ErrorAction SilentlyContinue)) { Write-Error "Go no instalado. https://go.dev/dl/"; exit 1 }
if (!(Get-Command node  -ErrorAction SilentlyContinue)) { Write-Error "Node.js no instalado. https://nodejs.org"; exit 1 }
if (!(Get-Command docker -ErrorAction SilentlyContinue)) { Write-Error "Docker no instalado. https://docker.com"; exit 1 }

Write-Host ("Go:     " + (go version)) -ForegroundColor Green
Write-Host ("Node:   " + (node --version)) -ForegroundColor Green

if (!(Test-Path ".env")) {
    Copy-Item ".env.example" ".env"
    Write-Host ".env creado - edita con tus API keys" -ForegroundColor Yellow
}

Write-Host "Instalando Go modules..." -ForegroundColor Cyan
Push-Location core; go mod download; go mod tidy; Pop-Location

Write-Host "Instalando Node.js packages..." -ForegroundColor Cyan
Push-Location scraper; npm install; npx playwright install chromium --with-deps; Pop-Location

New-Item -ItemType Directory -Force -Path "cv/base","cv/adapted","exports","bin" | Out-Null

Write-Host "Levantando Qdrant..." -ForegroundColor Cyan
docker compose up qdrant -d
Start-Sleep 4

Write-Host ""
Write-Host "Setup completo. Proximos pasos:" -ForegroundColor Green
Write-Host "  1. Edita .env con tu API key"
Write-Host "  2. Push-Location scraper; npm run dev"
Write-Host "  3. Push-Location core; go run ./cmd/offertrack"
