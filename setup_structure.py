"""
OfferTrack M82 - Generador de estructura del proyecto
Funciona en Windows, Mac y Linux.
Uso: python setup_structure.py
"""

import os

PROJECT = "offertrack-m82"

# ─── Helper ───────────────────────────────────────────────────────────────────
def write(path, content=""):
    full = os.path.join(PROJECT, path.replace("/", os.sep))
    os.makedirs(os.path.dirname(full), exist_ok=True)
    with open(full, "w", encoding="utf-8") as f:
        f.write(content)

def mkdir(path):
    os.makedirs(os.path.join(PROJECT, path.replace("/", os.sep)), exist_ok=True)

# ─────────────────────────────────────────────────────────────────────────────

os.makedirs(PROJECT, exist_ok=True)
print("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
print("  OfferTrack M82 — Generando estructura del proyecto")
print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

# ─── .github/workflows/ci.yml ────────────────────────────────────────────────
write(".github/workflows/ci.yml", """\
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  core:
    name: "Go - lint and test"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.24.2"
      - name: Install dependencies
        run: cd core && go mod download
      - name: Lint
        run: cd core && go vet ./...
      - name: Test
        run: cd core && go test -race ./...

  scraper:
    name: "Node.js - lint and build"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: "22"
      - name: Install dependencies
        run: cd scraper && npm ci
      - name: Lint
        run: cd scraper && npm run lint
      - name: Build
        run: cd scraper && npm run build
""")

# ─── core/cmd/offertrack/main.go ─────────────────────────────────────────────
write("core/cmd/offertrack/main.go", """\
package main

import (
\t"fmt"
\t"os"

\t"github.com/yourusername/offertrack-m82/core/internal/cli"
)

func main() {
\tif err := cli.Execute(); err != nil {
\t\tfmt.Fprintln(os.Stderr, err)
\t\tos.Exit(1)
\t}
}
""")

# ─── core/internal/cli ───────────────────────────────────────────────────────
write("core/internal/cli/root.go", """\
package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
\tUse:   "offertrack",
\tShort: "OfferTrack M82 - Plataforma de busqueda de empleo asistida por IA",
}

func Execute() error {
\treturn rootCmd.Execute()
}

func init() {
\trootCmd.AddCommand(searchCmd)
\trootCmd.AddCommand(analyzeCmd)
\trootCmd.AddCommand(adaptCmd)
\trootCmd.AddCommand(listCmd)
\trootCmd.AddCommand(configCmd)
}
""")

write("core/internal/cli/search.go", """\
package cli

import "github.com/spf13/cobra"

var searchCmd = &cobra.Command{
\tUse:   "search",
\tShort: "Buscar vacantes en portales de empleo",
\tRunE: func(cmd *cobra.Command, args []string) error {
\t\t// TODO: implementar busqueda
\t\treturn nil
\t},
}

func init() {
\tsearchCmd.Flags().String("role", "", "Puesto o area profesional")
\tsearchCmd.Flags().Int("salary-min", 0, "Salario minimo mensual")
\tsearchCmd.Flags().String("modality", "", "Modalidad: remote | hybrid | onsite")
\tsearchCmd.Flags().String("location", "", "Ciudad o municipio base")
}
""")

write("core/internal/cli/analyze.go", """\
package cli

import "github.com/spf13/cobra"

var analyzeCmd = &cobra.Command{
\tUse:   "analyze <job-id>",
\tShort: "Analizar una vacante contra tu perfil",
\tArgs:  cobra.ExactArgs(1),
\tRunE: func(cmd *cobra.Command, args []string) error {
\t\t// TODO: implementar analisis
\t\treturn nil
\t},
}
""")

write("core/internal/cli/adapt.go", """\
package cli

import "github.com/spf13/cobra"

var adaptCmd = &cobra.Command{
\tUse:   "adapt-cv <job-id>",
\tShort: "Generar CV adaptado para una vacante",
\tArgs:  cobra.ExactArgs(1),
\tRunE: func(cmd *cobra.Command, args []string) error {
\t\t// TODO: implementar adaptacion de CV
\t\treturn nil
\t},
}

func init() {
\tadaptCmd.Flags().String("output", "./exports", "Directorio de salida")
}
""")

write("core/internal/cli/list.go", """\
package cli

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
\tUse:   "list",
\tShort: "Listar vacantes recopiladas",
\tRunE: func(cmd *cobra.Command, args []string) error {
\t\t// TODO: implementar listado
\t\treturn nil
\t},
}
""")

write("core/internal/cli/config.go", """\
package cli

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
\tUse:   "config",
\tShort: "Gestionar configuracion del usuario y proveedor IA",
\tRunE: func(cmd *cobra.Command, args []string) error {
\t\t// TODO: implementar configuracion
\t\treturn nil
\t},
}
""")

# ─── core/internal/tui ───────────────────────────────────────────────────────
write("core/internal/tui/model.go", """\
package tui

import tea "github.com/charmbracelet/bubbletea"

type Model struct{}

func (m Model) Init() tea.Cmd                           { return nil }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m Model) View() string                            { return "" }
""")

write("core/internal/tui/styles/theme.go", """\
package styles

import "github.com/charmbracelet/lipgloss"

var (
\tTitleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
\tSubtitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
\tHighlightStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF88"))
\tDimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
)
""")

write("core/internal/tui/views/dashboard.go",  "package views\n")
write("core/internal/tui/views/job_list.go",   "package views\n")
write("core/internal/tui/views/job_detail.go", "package views\n")
write("core/internal/tui/views/cv_preview.go", "package views\n")

# ─── core/internal/ai ────────────────────────────────────────────────────────
write("core/internal/ai/provider.go", """\
package ai

import "context"

// AIProvider es el contrato que todos los proveedores deben cumplir.
// El nucleo nunca importa un proveedor concreto, solo esta interfaz.
type AIProvider interface {
\tAnalyze(ctx context.Context, req AnalysisRequest) (*AnalysisResult, error)
\tAdaptCV(ctx context.Context, req AdaptRequest) (*AdaptResult, error)
\tSummarize(ctx context.Context, text string) (string, error)
\tName() string
}

type AnalysisRequest struct {
\tJobDescription string
\tUserProfile    string
\tUserCV         string
}

type AnalysisResult struct {
\tCompatibilityScore int
\tStrengths          []string
\tGaps               []string
\tSalaryEstimate     string
\tRecommendation     string // "apply" | "consider" | "discard"
\tRawAnalysis        string
}

type AdaptRequest struct {
\tJobDescription string
\tBaseCV         string
\tUserProfile    string
}

type AdaptResult struct {
\tAdaptedCV     string
\tChanges       []string
\tKeywordsAdded []string
}
""")

write("core/internal/ai/factory.go", """\
package ai

import (
\t"fmt"

\t"github.com/yourusername/offertrack-m82/core/internal/ai/providers"
\t"github.com/yourusername/offertrack-m82/core/internal/config"
)

// NewProvider crea el proveedor correcto segun la configuracion.
// Agregar uno nuevo = archivo en providers/ + un case aqui.
func NewProvider(cfg *config.Config) (AIProvider, error) {
\tswitch cfg.AI.Provider {
\tcase "claude":
\t\treturn providers.NewClaudeProvider(cfg.AI.APIKey, cfg.AI.Model)
\tcase "gemini":
\t\treturn providers.NewGeminiProvider(cfg.AI.APIKey, cfg.AI.Model)
\tcase "groq":
\t\treturn providers.NewGroqProvider(cfg.AI.APIKey, cfg.AI.Model)
\tcase "openrouter":
\t\treturn providers.NewOpenRouterProvider(cfg.AI.APIKey, cfg.AI.Model)
\tdefault:
\t\treturn nil, fmt.Errorf("proveedor IA no soportado: %s", cfg.AI.Provider)
\t}
}
""")

write("core/internal/ai/prompt_builder.go", """\
package ai

import "fmt"

func BuildAnalysisPrompt(job, profile, cv string) string {
\treturn fmt.Sprintf(`Eres un experto en reclutamiento. Analiza esta vacante contra el perfil del candidato.

VACANTE:
%s

PERFIL DEL CANDIDATO:
%s

CV BASE:
%s

Responde SOLO en JSON con este esquema:
{
  "compatibility_score": <0-100>,
  "strengths": ["..."],
  "gaps": ["..."],
  "salary_estimate": "...",
  "recommendation": "apply|consider|discard",
  "raw_analysis": "..."
}`, job, profile, cv)
}

func BuildAdaptPrompt(job, cv, profile string) string {
\treturn fmt.Sprintf(`Eres un experto en redaccion de CVs. Adapta el CV a la vacante.

REGLAS:
- No inventes experiencias ni habilidades inexistentes
- Toda informacion debe estar respaldada en el CV base
- Preserva la linea de tiempo laboral
- Prioriza palabras clave de la vacante

VACANTE: %s
CV BASE: %s
PERFIL: %s

Responde SOLO en JSON:
{"adapted_cv":"<Markdown>","changes":["..."],"keywords_added":["..."]}`, job, cv, profile)
}
""")

# ─── providers ───────────────────────────────────────────────────────────────
for name, key_env in [
    ("claude",      "ANTHROPIC_API_KEY"),
    ("gemini",      "GEMINI_API_KEY"),
    ("groq",        "GROQ_API_KEY"),
    ("openrouter",  "OPENROUTER_API_KEY"),
]:
    struct = name.capitalize() + "Provider"
    write(f"core/internal/ai/providers/{name}.go", f"""\
package providers

import (
\t"context"
\t"fmt"

\t"github.com/yourusername/offertrack-m82/core/internal/ai"
)

type {struct} struct{{ apiKey, model string }}

func New{struct}(apiKey, model string) (*{struct}, error) {{
\tif apiKey == "" {{
\t\treturn nil, fmt.Errorf("{key_env} no configurada")
\t}}
\treturn &{struct}{{apiKey: apiKey, model: model}}, nil
}}

func (p *{struct}) Name() string {{ return "{name}" }}

func (p *{struct}) Analyze(ctx context.Context, req ai.AnalysisRequest) (*ai.AnalysisResult, error) {{
\treturn nil, nil // TODO: implementar
}}

func (p *{struct}) AdaptCV(ctx context.Context, req ai.AdaptRequest) (*ai.AdaptResult, error) {{
\treturn nil, nil // TODO: implementar
}}

func (p *{struct}) Summarize(ctx context.Context, text string) (string, error) {{
\treturn "", nil // TODO: implementar
}}
""")

# ─── core/internal/db ────────────────────────────────────────────────────────
write("core/internal/db/qdrant.go", """\
package db

import "github.com/qdrant/go-client/qdrant"

type QdrantClient struct {
\tclient *qdrant.Client
}

func NewQdrantClient(host string, port int) (*QdrantClient, error) {
\tclient, err := qdrant.NewClient(&qdrant.Config{
\t\tHost: host,
\t\tPort: port,
\t})
\tif err != nil {
\t\treturn nil, err
\t}
\treturn &QdrantClient{client: client}, nil
}
""")

write("core/internal/db/collections.go", """\
package db

const (
\tCollectionJobs    = "jobs"
\tCollectionProfile = "profile"
\tCollectionCVs     = "cv_versions"
\tCollectionMemory  = "claude_memory"
\tEmbedDimensions   = 384
)
""")

write("core/internal/db/queries.go", "package db\n")

# ─── core/internal/domain ────────────────────────────────────────────────────
write("core/internal/domain/job.go", """\
package domain

import "time"

type Job struct {
\tID          string    `json:"id"`
\tTitle       string    `json:"title"`
\tCompany     string    `json:"company"`
\tDescription string    `json:"description"`
\tSalary      string    `json:"salary"`
\tModality    string    `json:"modality"`
\tLocation    string    `json:"location"`
\tPortal      string    `json:"portal"`
\tURL         string    `json:"url"`
\tScrapedAt   time.Time `json:"scraped_at"`
\tCompatScore int       `json:"compat_score"`
\tStatus      string    `json:"status"`
}
""")

write("core/internal/domain/cv.go", """\
package domain

import "time"

type CV struct {
\tID          string    `json:"id"`
\tJobID       string    `json:"job_id"`
\tBaseContent string    `json:"base_content"`
\tAdapted     string    `json:"adapted"`
\tCreatedAt   time.Time `json:"created_at"`
}
""")

write("core/internal/domain/profile.go", """\
package domain

type Profile struct {
\tName       string   `json:"name"`
\tEmail      string   `json:"email"`
\tSkills     []string `json:"skills"`
\tExperience string   `json:"experience"`
\tEducation  string   `json:"education"`
\tSalaryMin  int      `json:"salary_min"`
\tModalities []string `json:"modalities"`
\tLocations  []string `json:"locations"`
}
""")

write("core/internal/domain/analysis.go", """\
package domain

type Analysis struct {
\tJobID              string   `json:"job_id"`
\tCompatibilityScore int      `json:"compatibility_score"`
\tStrengths          []string `json:"strengths"`
\tGaps               []string `json:"gaps"`
\tSalaryEstimate     string   `json:"salary_estimate"`
\tRecommendation     string   `json:"recommendation"`
\tRawAnalysis        string   `json:"raw_analysis"`
}
""")

# ─── core/internal/services ──────────────────────────────────────────────────
for svc in ["analyzer", "cv_adapter", "exporter", "orchestrator"]:
    write(f"core/internal/services/{svc}.go", "package services\n")

# ─── core/internal/config ────────────────────────────────────────────────────
write("core/internal/config/config.go", """\
package config

import "github.com/spf13/viper"

type Config struct {
\tApp     AppConfig
\tAI      AIConfig
\tQdrant  QdrantConfig
\tScraper ScraperConfig
}

type AppConfig   struct{ Name, Version, Env string }
type AIConfig    struct{ Provider, Model, APIKey string; Timeout int }
type QdrantConfig struct{ Host string; Port int }
type ScraperConfig struct{ BaseURL string; Timeout int }

func Load() (*Config, error) {
\tviper.SetConfigName("app")
\tviper.SetConfigType("yaml")
\tviper.AddConfigPath("./config")
\tviper.AutomaticEnv()
\tif err := viper.ReadInConfig(); err != nil {
\t\treturn nil, err
\t}
\tcfg := &Config{}
\treturn cfg, viper.Unmarshal(cfg)
}
""")

write("core/internal/config/defaults.go", """\
package config

import "github.com/spf13/viper"

func SetDefaults() {
\tviper.SetDefault("app.name", "OfferTrack M82")
\tviper.SetDefault("app.version", "1.0.0")
\tviper.SetDefault("app.env", "development")
\tviper.SetDefault("ai.provider", "gemini")
\tviper.SetDefault("ai.model", "gemini-2.0-flash")
\tviper.SetDefault("ai.timeout", 30)
\tviper.SetDefault("qdrant.host", "localhost")
\tviper.SetDefault("qdrant.port", 6334)
\tviper.SetDefault("scraper.base_url", "http://localhost:3001")
\tviper.SetDefault("scraper.timeout", 60)
}
""")

# ─── core/go.mod + Makefile ──────────────────────────────────────────────────
write("core/go.mod", """\
module github.com/yourusername/offertrack-m82/core

go 1.24.2

require (
\tgithub.com/charmbracelet/bubbletea v1.3.4
\tgithub.com/charmbracelet/lipgloss  v1.1.0
\tgithub.com/joho/godotenv           v1.5.1
\tgithub.com/qdrant/go-client        v1.14.0
\tgithub.com/spf13/cobra             v1.9.1
\tgithub.com/spf13/viper             v1.20.0
)
""")

write("core/Makefile", """\
.PHONY: build run dev test lint clean release

build:
\tgo build -o ../bin/offertrack ./cmd/offertrack

dev:
\tgo run -race ./cmd/offertrack

run: build
\t../bin/offertrack

test:
\tgo test -v -race ./...

coverage:
\tgo test -coverprofile=coverage.out ./...
\tgo tool cover -html=coverage.out

lint:
\tgo vet ./...

clean:
\trm -f ../bin/offertrack coverage.out

release:
\tCGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/offertrack ./cmd/offertrack
""")

# ─── scraper/src ─────────────────────────────────────────────────────────────
write("scraper/src/server.ts", """\
import express from "express";
import dotenv from "dotenv";
import { scrapeRoutes } from "./routes/scrape.routes";
import { embedRoutes } from "./routes/embed.routes";
import { healthRoutes } from "./routes/health.routes";
import { errorHandler } from "./middleware/error-handler";

dotenv.config();

const app = express();
const PORT = process.env.SCRAPER_PORT || 3001;

app.use(express.json({ limit: "10mb" }));
app.use("/scrape", scrapeRoutes);
app.use("/embed", embedRoutes);
app.use("/health", healthRoutes);
app.use(errorHandler);

app.listen(PORT, () => {
  console.log(`OfferTrack Scraper running on http://localhost:${PORT}`);
});
""")

write("scraper/src/portals/base.portal.ts", """\
export interface JobRaw {
  title: string;
  company: string;
  description: string;
  salary?: string;
  modality?: string;
  location?: string;
  url: string;
  portal: string;
}

export interface SearchParams {
  role: string;
  location?: string;
  salaryMin?: number;
  modality?: string;
}

export abstract class BasePortal {
  abstract name: string;
  abstract scrape(params: SearchParams): Promise<JobRaw[]>;

  protected randomDelay(): Promise<void> {
    const min = parseInt(process.env.SCRAPER_DELAY_MIN_MS || "1500");
    const max = parseInt(process.env.SCRAPER_DELAY_MAX_MS || "4000");
    const ms = Math.floor(Math.random() * (max - min + 1)) + min;
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
}
""")

for portal, comment in [
    ("occ",           "OCC Mundial"),
    ("computrabajo",  "Computrabajo"),
    ("indeed",        "Indeed Mexico"),
]:
    cls = portal.capitalize() + "Portal"
    write(f"scraper/src/portals/{portal}.portal.ts", f"""\
import {{ BasePortal, JobRaw, SearchParams }} from "./base.portal";

export class {cls} extends BasePortal {{
  name = "{portal}";

  async scrape(params: SearchParams): Promise<JobRaw[]> {{
    // TODO: implementar scraping de {comment}
    return [];
  }}
}}
""")

write("scraper/src/embeddings/embed.service.ts", """\
import { EmbeddingModel, FlagEmbedding } from "fastembed";

let model: FlagEmbedding | null = null;

async function getModel(): Promise<FlagEmbedding> {
  if (!model) {
    model = await FlagEmbedding.init({ model: EmbeddingModel.BGEBaseENV15 });
  }
  return model;
}

export async function embedTexts(texts: string[]): Promise<number[][]> {
  const m = await getModel();
  const results: number[][] = [];
  for await (const batch of m.embed(texts, 32)) {
    results.push(...batch);
  }
  return results;
}

export async function embedQuery(text: string): Promise<number[]> {
  const m = await getModel();
  return m.queryEmbed(text);
}
""")

write("scraper/src/embeddings/models.ts", """\
export const EMBED_MODEL = "BAAI/bge-small-en-v1.5";
export const EMBED_DIMENSIONS = 384;
""")

write("scraper/src/routes/scrape.routes.ts", """\
import { Router } from "express";

export const scrapeRoutes = Router();

scrapeRoutes.post("/", async (req, res) => {
  // TODO: implementar scraping
  res.json({ jobs: [] });
});
""")

write("scraper/src/routes/embed.routes.ts", """\
import { Router } from "express";
import { embedTexts, embedQuery } from "../embeddings/embed.service";

export const embedRoutes = Router();

embedRoutes.post("/texts", async (req, res) => {
  const { texts } = req.body;
  const embeddings = await embedTexts(texts);
  res.json({ embeddings });
});

embedRoutes.post("/query", async (req, res) => {
  const { text } = req.body;
  const embedding = await embedQuery(text);
  res.json({ embedding });
});
""")

write("scraper/src/routes/health.routes.ts", """\
import { Router } from "express";

export const healthRoutes = Router();

healthRoutes.get("/", (_req, res) => {
  res.json({ status: "ok", service: "offertrack-scraper" });
});
""")

write("scraper/src/middleware/rate-limiter.ts", """\
import { Request, Response, NextFunction } from "express";

const requests = new Map<string, number>();

export function rateLimiter(req: Request, res: Response, next: NextFunction) {
  const key = req.ip || "local";
  const now = Date.now();
  const last = requests.get(key) || 0;
  if (now - last < 500) {
    return res.status(429).json({ error: "Too many requests" });
  }
  requests.set(key, now);
  next();
}
""")

write("scraper/src/middleware/error-handler.ts", """\
import { Request, Response, NextFunction } from "express";

export function errorHandler(
  err: Error,
  _req: Request,
  res: Response,
  _next: NextFunction
) {
  console.error(err.stack);
  res.status(500).json({ error: err.message });
}
""")

write("scraper/src/utils/user-agent.ts", """\
const USER_AGENTS = [
  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124.0.0.0 Safari/537.36",
  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/124.0.0.0 Safari/537.36",
  "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/124.0.0.0 Safari/537.36",
];

export function randomUserAgent(): string {
  return USER_AGENTS[Math.floor(Math.random() * USER_AGENTS.length)];
}
""")

write("scraper/src/utils/delay.ts", """\
export function randomDelay(minMs = 1500, maxMs = 4000): Promise<void> {
  const ms = Math.floor(Math.random() * (maxMs - minMs + 1)) + minMs;
  return new Promise((resolve) => setTimeout(resolve, ms));
}
""")

write("scraper/src/utils/normalizer.ts", """\
export function normalizeText(text: string): string {
  return text.trim().replace(/\\s+/g, " ").toLowerCase();
}

export function normalizeJob(raw: Record<string, unknown>): Record<string, unknown> {
  return {
    ...raw,
    title: normalizeText(String(raw.title || "")),
    company: normalizeText(String(raw.company || "")),
  };
}
""")

write("scraper/package.json", """\
{
  "name": "offertrack-scraper",
  "version": "1.0.0",
  "description": "OfferTrack M82 - Modulo de recoleccion de vacantes",
  "main": "dist/server.js",
  "scripts": {
    "dev": "tsx watch src/server.ts",
    "build": "tsc",
    "start": "node dist/server.js",
    "lint": "eslint src --ext .ts",
    "test": "jest"
  },
  "dependencies": {
    "cheerio": "^1.0.0",
    "dotenv": "^16.0.0",
    "express": "^4.21.0",
    "fastembed": "^1.14.0",
    "playwright": "^1.51.0"
  },
  "devDependencies": {
    "@types/express": "^4.17.21",
    "@types/jest": "^29.0.0",
    "@types/node": "^22.0.0",
    "@typescript-eslint/eslint-plugin": "^8.0.0",
    "@typescript-eslint/parser": "^8.0.0",
    "eslint": "^9.0.0",
    "jest": "^29.0.0",
    "prettier": "^3.0.0",
    "tsx": "^4.0.0",
    "typescript": "^5.8.0"
  },
  "engines": {
    "node": ">=22.0.0"
  }
}
""")

write("scraper/tsconfig.json", """\
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
    "sourceMap": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
""")

write("scraper/.env.example", """\
SCRAPER_PORT=3001
SCRAPER_DELAY_MIN_MS=1500
SCRAPER_DELAY_MAX_MS=4000
""")

# ─── docker/ ─────────────────────────────────────────────────────────────────
write("docker/Dockerfile.core", """\
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
""")

write("docker/Dockerfile.scraper", """\
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
""")

write("docker/qdrant/config.yaml", """\
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
    on_disk: false
  optimizer:
    default_segment_number: 2
  replication_factor: 1
""")

# ─── config/ ─────────────────────────────────────────────────────────────────
write("config/app.yaml", """\
app:
  name: "OfferTrack M82"
  version: "1.0.0"
  env: development

ai:
  provider: gemini
  model: gemini-2.0-flash
  timeout_seconds: 30
  max_retries: 3

qdrant:
  host: localhost
  port: 6334
  collections:
    jobs: jobs
    profile: profile
    cvs: cv_versions
    memory: claude_memory

scraper:
  base_url: http://localhost:3001
  timeout_seconds: 60

embedding:
  model: BAAI/bge-small-en-v1.5
  dimensions: 384

export:
  output_dir: ./exports
  formats:
    - markdown
    - pdf
    - docx
""")

write("config/providers.yaml", """\
providers:
  claude:
    base_url: https://api.anthropic.com
    default_model: claude-sonnet-4-5
    env_key: ANTHROPIC_API_KEY

  gemini:
    base_url: https://generativelanguage.googleapis.com
    default_model: gemini-2.0-flash
    env_key: GEMINI_API_KEY

  groq:
    base_url: https://api.groq.com/openai/v1
    default_model: llama-3.3-70b-versatile
    env_key: GROQ_API_KEY

  openrouter:
    base_url: https://openrouter.ai/api/v1
    default_model: meta-llama/llama-3.3-70b-instruct
    env_key: OPENROUTER_API_KEY
""")

# ─── scripts/ ────────────────────────────────────────────────────────────────
write("scripts/setup.ps1", """\
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
""")

write("scripts/reset-db.ps1", """\
$confirm = Read-Host "Esto eliminara todos los datos de Qdrant. Continuar? (s/N)"
if ($confirm -eq "s" -or $confirm -eq "S") {
    docker compose down -v
    docker compose up qdrant -d
    Write-Host "Qdrant reiniciado." -ForegroundColor Green
} else {
    Write-Host "Cancelado."
}
""")

write("scripts/seed-profile.ps1", """\
Write-Host "Cargando perfil base en Qdrant..." -ForegroundColor Cyan
# TODO: implementar seed del perfil de usuario
Write-Host "Listo." -ForegroundColor Green
""")

# ─── cv/ ─────────────────────────────────────────────────────────────────────
mkdir("cv/adapted")
write("cv/base/cv_base.md", """\
# Tu Nombre

**Email:** tu@email.com | **LinkedIn:** linkedin.com/in/tuperfil | **GitHub:** github.com/tuusuario

---

## Resumen Profesional

Describe brevemente tu perfil profesional aqui.

---

## Experiencia

### Puesto - Empresa - Fecha inicio / Fecha fin

- Logro o responsabilidad 1
- Logro o responsabilidad 2

---

## Habilidades

**Lenguajes:** Go, TypeScript, Python
**Herramientas:** Docker, Git, Linux
**Bases de datos:** PostgreSQL, Redis

---

## Educacion

**Grado - Institucion - Anio de graduacion**
""")

# ─── exports/ + docs/ + bin/ ─────────────────────────────────────────────────
mkdir("exports")
mkdir("bin")
write("docs/architecture.md", "")
write("docs/providers.md",    "")
write("docs/api-scraper.md",  "")

# ─── docker-compose.yml ──────────────────────────────────────────────────────
write("docker-compose.yml", """\
version: "3.9"

services:

  qdrant:
    image: qdrant/qdrant:v1.16.0
    container_name: offertrack-qdrant
    ports:
      - "6333:6333"
      - "6334:6334"
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
      - ./scraper/src:/app/src
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
""")

# ─── .env.example ────────────────────────────────────────────────────────────
write(".env.example", """\
# Proveedor IA activo
AI_PROVIDER=gemini
AI_MODEL=gemini-2.0-flash

# API Keys (solo la del proveedor activo es requerida)
ANTHROPIC_API_KEY=sk-ant-xxxxxxxxxxxxxxxxxxxx
GEMINI_API_KEY=AIzaSyxxxxxxxxxxxxxxxxxxxxxxxxx
GROQ_API_KEY=gsk_xxxxxxxxxxxxxxxxxxxxxxxxxx
OPENROUTER_API_KEY=sk-or-xxxxxxxxxxxxxxxxxxxx

# Qdrant
QDRANT_HOST=localhost
QDRANT_PORT=6334
QDRANT_COLLECTION_JOBS=jobs
QDRANT_COLLECTION_PROFILE=profile
QDRANT_COLLECTION_CVS=cv_versions
QDRANT_COLLECTION_MEMORY=claude_memory

# Scraper
SCRAPER_PORT=3001
SCRAPER_BASE_URL=http://localhost:3001
SCRAPER_DELAY_MIN_MS=1500
SCRAPER_DELAY_MAX_MS=4000

# Embedding
EMBED_MODEL=BAAI/bge-small-en-v1.5
EMBED_DIMENSIONS=384

# App
LOG_LEVEL=info
ENV=development
""")

# ─── .gitignore ──────────────────────────────────────────────────────────────
write(".gitignore", """\
# Entorno
.env
.env.local

# Go
bin/
core/coverage.out
core/go.sum

# Node.js
scraper/node_modules/
scraper/dist/

# Datos de usuario
cv/base/
cv/adapted/
exports/

# Qdrant local
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
""")

# ─── Resumen ─────────────────────────────────────────────────────────────────
total = sum(len(files) for _, _, files in os.walk(PROJECT))
print(f"\n  Archivos generados: {total}")
print(f"  Directorio:         {os.path.abspath(PROJECT)}")
print("\n  Proximos pasos:")
print(f"    cd {PROJECT}")
print("    python --version  (requiere Python 3.x, ya lo tienes)")
print("    Edita .env con tu API key")
print("    Ejecuta scripts\\setup.ps1 para instalar dependencias")
print("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
