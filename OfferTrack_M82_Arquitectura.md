# OfferTrack M82 — Arquitectura Completa del Proyecto

> Plataforma de búsqueda de empleo asistida por IA  
> Versión de arquitectura: 1.0 | Abril 2026

---

## 1. Árbol de directorios

```
offertrack-m82/
│
├── .github/
│   └── workflows/
│       └── ci.yml                    # Pipeline CI: lint, test, build
│
├── core/                             # ── Módulo Go (Orquestador principal)
│   ├── cmd/
│   │   └── offertrack/
│   │       └── main.go               # Punto de entrada del binario
│   │
│   ├── internal/
│   │   ├── cli/
│   │   │   ├── root.go               # Comando raíz Cobra
│   │   │   ├── search.go             # Comando: offertrack search
│   │   │   ├── analyze.go            # Comando: offertrack analyze <id>
│   │   │   ├── adapt.go              # Comando: offertrack adapt-cv <id>
│   │   │   ├── list.go               # Comando: offertrack list
│   │   │   └── config.go             # Comando: offertrack config
│   │   │
│   │   ├── tui/
│   │   │   ├── model.go              # Modelo principal Bubble Tea
│   │   │   ├── views/
│   │   │   │   ├── dashboard.go      # Vista: panel principal
│   │   │   │   ├── job_list.go       # Vista: lista de vacantes
│   │   │   │   ├── job_detail.go     # Vista: detalle de vacante
│   │   │   │   └── cv_preview.go     # Vista: previsualización de CV
│   │   │   └── styles/
│   │   │       └── theme.go          # Paleta de colores Lip Gloss
│   │   │
│   │   ├── ai/
│   │   │   ├── provider.go           # Interfaz AIProvider (contrato)
│   │   │   ├── factory.go            # Factory: crea proveedor según config
│   │   │   ├── prompt_builder.go     # Constructor de prompts reutilizables
│   │   │   └── providers/
│   │   │       ├── claude.go         # Implementación: Claude (Anthropic)
│   │   │       ├── gemini.go         # Implementación: Gemini (Google)
│   │   │       ├── groq.go           # Implementación: Groq (Llama/DeepSeek)
│   │   │       └── openrouter.go     # Implementación: OpenRouter
│   │   │
│   │   ├── db/
│   │   │   ├── qdrant.go             # Cliente Qdrant + helpers
│   │   │   ├── collections.go        # Definición de colecciones y schemas
│   │   │   └── queries.go            # Queries reutilizables
│   │   │
│   │   ├── domain/
│   │   │   ├── job.go                # Entidad: Vacante
│   │   │   ├── cv.go                 # Entidad: Currículum
│   │   │   ├── profile.go            # Entidad: Perfil de usuario
│   │   │   └── analysis.go           # Entidad: Resultado de análisis
│   │   │
│   │   ├── services/
│   │   │   ├── analyzer.go           # Servicio: análisis de vacantes
│   │   │   ├── cv_adapter.go         # Servicio: adaptación de CV
│   │   │   ├── exporter.go           # Servicio: exportación PDF/DOCX
│   │   │   └── orchestrator.go       # Servicio: orquestación del flujo completo
│   │   │
│   │   └── config/
│   │       ├── config.go             # Carga y validación de configuración
│   │       └── defaults.go           # Valores por defecto
│   │
│   ├── go.mod                        # Módulo Go y dependencias
│   ├── go.sum
│   └── Makefile                      # Comandos de desarrollo
│
├── scraper/                          # ── Módulo Node.js (Recolector de vacantes)
│   ├── src/
│   │   ├── server.ts                 # Servidor HTTP Express (API local)
│   │   │
│   │   ├── portals/
│   │   │   ├── base.portal.ts        # Clase base: interfaz común de scraping
│   │   │   ├── occ.portal.ts         # Scraper: OCC Mundial
│   │   │   ├── computrabajo.portal.ts # Scraper: Computrabajo
│   │   │   └── indeed.portal.ts      # Scraper: Indeed México
│   │   │
│   │   ├── embeddings/
│   │   │   ├── embed.service.ts      # Servicio: generación de embeddings
│   │   │   └── models.ts             # Configuración de modelos fastembed
│   │   │
│   │   ├── routes/
│   │   │   ├── scrape.routes.ts      # POST /scrape
│   │   │   ├── embed.routes.ts       # POST /embed
│   │   │   └── health.routes.ts      # GET /health
│   │   │
│   │   ├── middleware/
│   │   │   ├── rate-limiter.ts       # Control de frecuencia de requests
│   │   │   └── error-handler.ts      # Manejo centralizado de errores
│   │   │
│   │   └── utils/
│   │       ├── user-agent.ts         # Rotación de User-Agent
│   │       ├── delay.ts              # Delays aleatorios entre requests
│   │       └── normalizer.ts         # Normalización de datos scrapeados
│   │
│   ├── package.json
│   ├── tsconfig.json
│   └── .env.example
│
├── docker/                           # ── Contenedores
│   ├── Dockerfile.core               # Imagen Go
│   ├── Dockerfile.scraper            # Imagen Node.js
│   └── qdrant/
│       └── config.yaml               # Configuración Qdrant
│
├── config/                           # ── Configuración global
│   ├── app.yaml                      # Config principal de la app
│   └── providers.yaml                # Config de proveedores IA
│
├── scripts/                          # ── Scripts de utilidad
│   ├── setup.sh                      # Setup completo del entorno
│   ├── seed-profile.sh               # Carga perfil inicial de usuario
│   └── reset-db.sh                   # Limpia y reinicia Qdrant
│
├── cv/                               # ── CVs del usuario (gitignored)
│   ├── base/
│   │   └── cv_base.md                # CV base en Markdown
│   └── adapted/                      # CVs generados por la app
│
├── exports/                          # ── Archivos exportados (gitignored)
│
├── docs/                             # ── Documentación técnica
│   ├── architecture.md
│   ├── providers.md
│   └── api-scraper.md
│
├── .env                              # Variables de entorno (gitignored)
├── .env.example                      # Plantilla de variables
├── .gitignore
├── docker-compose.yml                # Orquestación completa
└── README.md
```

---

## 2. Stack tecnológico

### Core (Go)

| Tecnología | Versión | Rol |
|---|---|---|
| Go | 1.24.2 | Lenguaje principal — binario único, rendimiento nativo |
| Cobra | v1.9.1 | Framework CLI — gestión de comandos y flags |
| Bubble Tea | v1.3.4 | TUI interactiva — vistas y navegación en terminal |
| Lip Gloss | v1.1.0 | Estilos visuales para la TUI |
| go-client (Qdrant) | v1.14.0 | Cliente gRPC oficial para Qdrant |
| viper | v1.20.0 | Gestión de configuración (YAML + env vars) |
| godotenv | v1.5.1 | Carga de archivos .env |

### Scraper (Node.js)

| Tecnología | Versión | Rol |
|---|---|---|
| Node.js | 22 LTS | Runtime — ecosistema maduro para scraping |
| TypeScript | 5.8.x | Tipado estático — previene errores en runtime |
| Playwright | 1.51.x | Scraping de portales con JavaScript dinámico |
| Cheerio | 1.0.0 | Parsing HTML estático — rápido y liviano |
| Express | 4.21.x | Servidor HTTP local — API que consume Go |
| fastembed | 1.14.x | Embeddings locales sin API externa |
| tsx | 4.x | Ejecución directa de TypeScript en desarrollo |

### Base de datos

| Tecnología | Versión | Rol |
|---|---|---|
| Qdrant | v1.16.0 | DB única — vectores + payload estructurado |
| Modelo embedding | BAAI/bge-small-en-v1.5 | 384 dimensiones, equilibrio velocidad/precisión |

### IA — Proveedores (intercambiables)

| Proveedor | Modelo recomendado | Tier |
|---|---|---|
| Claude (Anthropic) | claude-sonnet-4-5 | Pago por uso |
| Gemini (Google) | gemini-2.0-flash | Gratuito con límites |
| Groq | llama-3.3-70b-versatile | Gratuito con límites |
| OpenRouter | (configurable) | Gratuito / Pago |

### Exportación

| Tecnología | Versión | Rol |
|---|---|---|
| Pandoc | 3.6.x | Conversión MD → PDF / DOCX |

---

## 3. Entorno de desarrollo

### Editor recomendado: VS Code

**Descarga:** https://code.visualstudio.com/

**Extensiones requeridas** (instalar en orden):

```
Go                    → golang.go
ESLint                → dbaeumer.vscode-eslint
Prettier              → esbenp.prettier-vscode
Docker                → ms-azuretools.vscode-docker
YAML                  → redhat.vscode-yaml
DotENV                → mikestead.dotenv
```

**`.vscode/settings.json` recomendado:**

```json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  },
  "go.toolsManagement.autoUpdate": true,
  "go.lintTool": "golangci-lint",
  "typescript.preferences.importModuleSpecifier": "relative",
  "files.exclude": {
    "**/node_modules": true,
    "**/.git": true,
    "**/exports": true
  }
}
```

### Herramientas de sistema requeridas

| Herramienta | Versión | Instalación |
|---|---|---|
| Go | 1.24.2 | https://go.dev/dl/ |
| Node.js | 22 LTS | https://nodejs.org o `nvm install 22` |
| Docker Desktop | Latest | https://www.docker.com/products/docker-desktop/ |
| Pandoc | 3.6.x | https://pandoc.org/installing.html |
| Git | 2.x | https://git-scm.com/downloads |

**Instalar nvm (recomendado para Node.js):**

```bash
# macOS / Linux
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
nvm install 22
nvm use 22

# Windows
# Descargar nvm-windows: https://github.com/coreybutler/nvm-windows/releases
```

---

## 4. Archivos de configuración base

### `core/go.mod`

```go
module github.com/yourusername/offertrack-m82/core

go 1.24.2

require (
    github.com/spf13/cobra       v1.9.1
    github.com/charmbracelet/bubbletea v1.3.4
    github.com/charmbracelet/lipgloss  v1.1.0
    github.com/qdrant/go-client        v1.14.0
    github.com/spf13/viper             v1.20.0
    github.com/joho/godotenv           v1.5.1
)
```

### `scraper/package.json`

```json
{
  "name": "offertrack-scraper",
  "version": "1.0.0",
  "description": "OfferTrack M82 — Módulo de recolección de vacantes",
  "main": "src/server.ts",
  "scripts": {
    "dev": "tsx watch src/server.ts",
    "build": "tsc",
    "start": "node dist/server.js",
    "lint": "eslint src --ext .ts",
    "test": "jest"
  },
  "dependencies": {
    "express": "^4.21.0",
    "playwright": "^1.51.0",
    "cheerio": "^1.0.0",
    "fastembed": "^1.14.0",
    "dotenv": "^16.0.0"
  },
  "devDependencies": {
    "typescript": "^5.8.0",
    "@types/express": "^4.17.21",
    "@types/node": "^22.0.0",
    "tsx": "^4.0.0",
    "eslint": "^9.0.0",
    "@typescript-eslint/parser": "^8.0.0",
    "@typescript-eslint/eslint-plugin": "^8.0.0",
    "prettier": "^3.0.0",
    "jest": "^29.0.0",
    "@types/jest": "^29.0.0"
  },
  "engines": {
    "node": ">=22.0.0"
  }
}
```

### `scraper/tsconfig.json`

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "lib": ["ES2022"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
```

---

## 5. Variables de entorno

### `.env.example`

```env
# ─── Proveedor IA activo ───────────────────────────────────────
AI_PROVIDER=gemini           # claude | gemini | groq | openrouter
AI_MODEL=gemini-2.0-flash    # Modelo específico del proveedor

# ─── API Keys (solo la del proveedor activo es requerida) ──────
ANTHROPIC_API_KEY=sk-ant-xxxxxxxxxxxxxxxxxxxx
GEMINI_API_KEY=AIzaSyxxxxxxxxxxxxxxxxxxxxxxxxx
GROQ_API_KEY=gsk_xxxxxxxxxxxxxxxxxxxxxxxxxx
OPENROUTER_API_KEY=sk-or-xxxxxxxxxxxxxxxxxxxx

# ─── Qdrant ───────────────────────────────────────────────────
QDRANT_HOST=localhost
QDRANT_PORT=6334             # gRPC (más rápido que REST 6333)
QDRANT_COLLECTION_JOBS=jobs
QDRANT_COLLECTION_PROFILE=profile
QDRANT_COLLECTION_CVS=cv_versions
QDRANT_COLLECTION_MEMORY=claude_memory

# ─── Scraper ──────────────────────────────────────────────────
SCRAPER_PORT=3001
SCRAPER_BASE_URL=http://localhost:3001
SCRAPER_DELAY_MIN_MS=1500    # Delay mínimo entre requests
SCRAPER_DELAY_MAX_MS=4000    # Delay máximo entre requests

# ─── Embedding ────────────────────────────────────────────────
EMBED_MODEL=BAAI/bge-small-en-v1.5
EMBED_DIMENSIONS=384

# ─── App ──────────────────────────────────────────────────────
LOG_LEVEL=info               # debug | info | warn | error
ENV=development              # development | production
```

### `config/app.yaml`

```yaml
app:
  name: "OfferTrack M82"
  version: "1.0.0"
  env: ${ENV}

ai:
  provider: ${AI_PROVIDER}
  model: ${AI_MODEL}
  timeout_seconds: 30
  max_retries: 3

qdrant:
  host: ${QDRANT_HOST}
  port: ${QDRANT_PORT}
  collections:
    jobs: ${QDRANT_COLLECTION_JOBS}
    profile: ${QDRANT_COLLECTION_PROFILE}
    cvs: ${QDRANT_COLLECTION_CVS}
    memory: ${QDRANT_COLLECTION_MEMORY}

scraper:
  base_url: ${SCRAPER_BASE_URL}
  timeout_seconds: 60

embedding:
  model: ${EMBED_MODEL}
  dimensions: ${EMBED_DIMENSIONS}

export:
  output_dir: ./exports
  formats:
    - markdown
    - pdf
    - docx
```

---

## 6. Docker

### `docker-compose.yml`

```yaml
version: "3.9"

services:

  qdrant:
    image: qdrant/qdrant:v1.16.0
    container_name: offertrack-qdrant
    ports:
      - "6333:6333"   # REST + Dashboard web
      - "6334:6334"   # gRPC (Go se conecta aquí)
    volumes:
      - qdrant_data:/qdrant/storage
      - ./docker/qdrant/config.yaml:/qdrant/config/production.yaml
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:6333/healthz"]
      interval: 10s
      timeout: 5s
      retries: 5

  scraper:
    build:
      context: ./scraper
      dockerfile: ../docker/Dockerfile.scraper
    container_name: offertrack-scraper
    ports:
      - "3001:3001"
    env_file:
      - .env
    volumes:
      - ./scraper/src:/app/src   # Hot reload en desarrollo
    depends_on:
      qdrant:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3001/health"]
      interval: 15s
      timeout: 5s
      retries: 3

volumes:
  qdrant_data:
    driver: local
```

### `docker/Dockerfile.scraper`

```dockerfile
FROM node:22-alpine AS base
WORKDIR /app

FROM base AS deps
COPY package*.json ./
RUN npm ci --only=production

FROM base AS builder
COPY package*.json ./
RUN npm ci
COPY . .
RUN npx playwright install chromium --with-deps
RUN npm run build

FROM base AS runner
ENV NODE_ENV=production
COPY --from=deps /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist
COPY --from=builder /root/.cache/ms-playwright /root/.cache/ms-playwright

EXPOSE 3001
CMD ["node", "dist/server.js"]
```

### `docker/Dockerfile.core`

```dockerfile
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o offertrack ./cmd/offertrack

FROM alpine:3.21 AS runner
RUN apk add --no-cache pandoc ca-certificates
WORKDIR /app
COPY --from=builder /app/offertrack .
COPY config/ ./config/

ENTRYPOINT ["./offertrack"]
```

### `docker/qdrant/config.yaml`

```yaml
storage:
  storage_path: /qdrant/storage

service:
  host: 0.0.0.0
  http_port: 6333
  grpc_port: 6334
  max_request_size_mb: 32

log_level: INFO

collection:
  vectors:
    on_disk: false          # En RAM para máxima velocidad (escala personal)
  optimizer:
    default_segment_number: 2
  replication_factor: 1
```

---

## 7. Interfaz AIProvider (código base)

### `core/internal/ai/provider.go`

```go
package ai

import "context"

// AIProvider es el contrato que todos los proveedores deben cumplir.
// El núcleo nunca importa un proveedor concreto — solo esta interfaz.
type AIProvider interface {
    // Analyze evalúa una vacante contra el perfil del usuario
    Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResult, error)

    // AdaptCV reescribe el CV base orientado a una vacante específica
    AdaptCV(ctx context.Context, req AdaptRequest) (*AdaptResult, error)

    // Summarize genera un resumen conciso de texto largo
    Summarize(ctx context.Context, text string) (string, error)

    // Name retorna el nombre del proveedor activo (para logs y UI)
    Name() string
}

type AnalysisRequest struct {
    JobDescription string
    UserProfile    string
    UserCV         string
}

type AnalysisResult struct {
    CompatibilityScore int      // 0-100
    Strengths          []string
    Gaps               []string
    SalaryEstimate     string
    Recommendation     string   // "apply" | "consider" | "discard"
    RawAnalysis        string
}

type AdaptRequest struct {
    JobDescription string
    BaseCV         string
    UserProfile    string
}

type AdaptResult struct {
    AdaptedCV   string // Markdown
    Changes     []string
    KeywordsAdded []string
}
```

### `core/internal/ai/factory.go`

```go
package ai

import (
    "fmt"
    "github.com/yourusername/offertrack-m82/core/internal/ai/providers"
    "github.com/yourusername/offertrack-m82/core/internal/config"
)

// NewProvider crea el proveedor correcto según la configuración.
// Agregar un proveedor nuevo = agregar un case aquí + crear el archivo en providers/
func NewProvider(cfg *config.Config) (AIProvider, error) {
    switch cfg.AI.Provider {
    case "claude":
        return providers.NewClaudeProvider(cfg.AI.APIKey, cfg.AI.Model)
    case "gemini":
        return providers.NewGeminiProvider(cfg.AI.APIKey, cfg.AI.Model)
    case "groq":
        return providers.NewGroqProvider(cfg.AI.APIKey, cfg.AI.Model)
    case "openrouter":
        return providers.NewOpenRouterProvider(cfg.AI.APIKey, cfg.AI.Model)
    default:
        return nil, fmt.Errorf("proveedor IA no soportado: %s", cfg.AI.Provider)
    }
}
```

---

## 8. Instalación de dependencias

### Setup completo en un solo paso

```bash
# Clonar el repositorio
git clone https://github.com/yourusername/offertrack-m82.git
cd offertrack-m82

# Dar permisos al script de setup
chmod +x scripts/setup.sh

# Ejecutar setup completo
./scripts/setup.sh
```

### `scripts/setup.sh`

```bash
#!/usr/bin/env bash
set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  OfferTrack M82 — Setup de entorno"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 1. Verificar dependencias del sistema
command -v go     >/dev/null 2>&1 || { echo "❌ Go no instalado. https://go.dev/dl/"; exit 1; }
command -v node   >/dev/null 2>&1 || { echo "❌ Node.js no instalado. https://nodejs.org"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "❌ Docker no instalado. https://docker.com"; exit 1; }

echo "✅ Go:     $(go version)"
echo "✅ Node:   $(node --version)"
echo "✅ Docker: $(docker --version)"

# 2. Copiar .env si no existe
if [ ! -f .env ]; then
    cp .env.example .env
    echo "📄 .env creado — edítalo con tus API keys antes de continuar"
fi

# 3. Instalar dependencias Go
echo ""
echo "📦 Instalando dependencias Go..."
cd core && go mod download && go mod tidy
cd ..

# 4. Instalar dependencias Node.js
echo ""
echo "📦 Instalando dependencias Node.js..."
cd scraper && npm install
echo "🌐 Instalando browsers Playwright..."
npx playwright install chromium --with-deps
cd ..

# 5. Crear directorios necesarios
mkdir -p cv/base cv/adapted exports

# 6. Levantar Qdrant
echo ""
echo "🚀 Levantando Qdrant..."
docker compose up qdrant -d
sleep 3

# 7. Verificar que Qdrant está listo
curl -sf http://localhost:6333/healthz > /dev/null && \
    echo "✅ Qdrant listo en http://localhost:6333/dashboard" || \
    echo "⚠️  Qdrant tardando — verifica con: docker compose logs qdrant"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Setup completo."
echo "  Edita .env con tu API key y ejecuta:"
echo "  make dev"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
```

### Comandos manuales (si prefieres paso a paso)

```bash
# ── Go ────────────────────────────────────
cd core
go mod download
go mod tidy
go build ./...          # Verificar que compila
cd ..

# ── Node.js ───────────────────────────────
cd scraper
npm install
npx playwright install chromium --with-deps
npm run build           # Verificar que compila TypeScript
cd ..

# ── Qdrant (Docker) ───────────────────────
docker compose up qdrant -d
```

---

## 9. Makefile (core)

```makefile
# core/Makefile

.PHONY: build run dev test lint clean

# Compilar binario
build:
	@echo "Building OfferTrack M82..."
	@go build -o ../bin/offertrack ./cmd/offertrack

# Ejecutar en modo desarrollo (con race detector)
dev:
	@go run -race ./cmd/offertrack

# Ejecutar el binario compilado
run: build
	@../bin/offertrack

# Tests
test:
	@go test -v -race ./...

# Tests con cobertura
coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# Linter
lint:
	@golangci-lint run ./...

# Limpiar artefactos
clean:
	@rm -f ../bin/offertrack coverage.out

# Generar binario para producción (optimizado)
release:
	@CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/offertrack ./cmd/offertrack
```

---

## 10. Ejecución del proyecto

### Desarrollo (flujo completo)

```bash
# Terminal 1 — Infraestructura (Qdrant)
docker compose up qdrant

# Terminal 2 — Scraper
cd scraper && npm run dev

# Terminal 3 — Core Go
cd core && make dev
```

### Un solo comando (producción local)

```bash
# Levanta Qdrant + Scraper juntos en background
docker compose up -d

# Compilar y ejecutar la CLI
cd core && make build
./bin/offertrack --help
```

### Comandos de la CLI

```bash
# Configurar perfil de usuario
offertrack config

# Buscar vacantes (con parámetros)
offertrack search \
  --role "Backend Developer" \
  --salary-min 20000 \
  --modality remote \
  --location "Monterrey, NL"

# Listar vacantes recolectadas
offertrack list

# Analizar una vacante específica
offertrack analyze <job-id>

# Adaptar CV a una vacante
offertrack adapt-cv <job-id> --output ./exports/cv_empresa.md

# Cambiar proveedor IA sin reiniciar
offertrack config set-provider gemini
```

### Verificar estado del sistema

```bash
# Estado de todos los servicios
docker compose ps

# Dashboard Qdrant (browser)
open http://localhost:6333/dashboard

# Health del scraper
curl http://localhost:3001/health

# Logs en tiempo real
docker compose logs -f
```

### Apagar el entorno

```bash
# Apagar servicios (datos persisten en volumen Docker)
docker compose down

# Apagar Y borrar datos (reset completo)
docker compose down -v
```

---

## 11. `.gitignore`

```gitignore
# Entorno
.env
.env.local

# Go
core/bin/
core/coverage.out

# Node.js
scraper/node_modules/
scraper/dist/

# Datos de usuario
cv/base/
cv/adapted/
exports/

# Qdrant (si se usa volumen local en vez de Docker)
qdrant_storage/

# IDEs
.vscode/
.idea/
*.swp

# OS
.DS_Store
Thumbs.db

# Logs
*.log
```

---

## 12. Checklist de inicio rápido

```
□ 1. Instalar Go 1.24.2         → https://go.dev/dl/
□ 2. Instalar Node.js 22 LTS    → https://nodejs.org
□ 3. Instalar Docker Desktop     → https://docker.com
□ 4. Clonar el repositorio       → git clone ...
□ 5. Ejecutar setup              → ./scripts/setup.sh
□ 6. Editar .env                 → agregar API key del proveedor elegido
□ 7. Levantar infraestructura    → docker compose up -d
□ 8. Iniciar scraper             → cd scraper && npm run dev
□ 9. Compilar CLI                → cd core && make build
□ 10. Configurar perfil          → offertrack config
□ 11. Primera búsqueda           → offertrack search --role "Backend Dev"
```

---

*OfferTrack M82 — Arquitectura v1.0 | Abril 2026*
