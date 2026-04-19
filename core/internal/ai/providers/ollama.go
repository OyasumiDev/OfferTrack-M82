// core/internal/ai/providers/ollama.go
package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/QUERTY/OfferTrack-M82/internal/domain"
	"github.com/QUERTY/OfferTrack-M82/internal/prompts"
)

const ollamaDefaultBase = "http://localhost:11434"

// OllamaProvider implementa AIProvider usando la API local de Ollama.
// No requiere API key — corre 100% offline.
type OllamaProvider struct {
	baseURL string
	model   string
}

// NewOllamaProvider crea y valida el proveedor Ollama.
// Verifica que Ollama esté corriendo y que el modelo solicitado esté descargado.
func NewOllamaProvider(baseURL, model string) (*OllamaProvider, error) {
	if baseURL == "" {
		baseURL = ollamaDefaultBase
	}
	if model == "" {
		model = "deepseek-r1"
	}

	models, err := ollamaListModels(baseURL)
	if err != nil {
		return nil, fmt.Errorf(
			"Ollama no está corriendo en %s → %w\n(Ábrelo desde el menú inicio o ejecuta: ollama serve)",
			baseURL, err,
		)
	}

	if !ollamaModelExists(models, model) {
		return nil, fmt.Errorf(
			"modelo %q no encontrado en Ollama\nDescárgalo con: ollama pull %s\nModelos disponibles: %s",
			model, model, strings.Join(models, ", "),
		)
	}

	return &OllamaProvider{baseURL: baseURL, model: model}, nil
}

func (p *OllamaProvider) Name() string { return "ollama/" + p.model }

// ── Métodos AIProvider ────────────────────────────────────────────────────────

func (p *OllamaProvider) Analyze(ctx context.Context, req domain.AnalysisRequest) (*domain.AnalysisResult, error) {
	sys, user := prompts.BuildAnalysisMessages(req)
	raw, err := p.chat(ctx, sys, user, "json", 120*time.Second)
	if err != nil {
		return nil, fmt.Errorf("ollama analyze: %w", err)
	}
	return prompts.ParseAnalysisResult(raw)
}

func (p *OllamaProvider) AdaptCV(ctx context.Context, req domain.AdaptRequest) (*domain.AdaptResult, error) {
	sys, user := prompts.BuildAdaptMessages(req)
	raw, err := p.chat(ctx, sys, user, "json", 180*time.Second)
	if err != nil {
		return nil, fmt.Errorf("ollama adapt_cv: %w", err)
	}
	return prompts.ParseAdaptResult(raw)
}

func (p *OllamaProvider) Summarize(ctx context.Context, text string) (string, error) {
	sys, user := prompts.BuildSummarizeMessages(text)
	raw, err := p.chat(ctx, sys, user, "", 60*time.Second)
	if err != nil {
		return "", fmt.Errorf("ollama summarize: %w", err)
	}
	return raw, nil
}

// ── Implementación HTTP ───────────────────────────────────────────────────────

type ollamaChatReq struct {
	Model    string   `json:"model"`
	Messages []oaiMsg `json:"messages"` // reutiliza oaiMsg de http_helper.go
	Stream   bool     `json:"stream"`
	Format   string   `json:"format,omitempty"`
}

type ollamaChatResp struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Error string `json:"error,omitempty"`
	Done  bool   `json:"done"`
}

// chat llama a /api/chat de Ollama con el sistema y mensaje del usuario.
// format = "json" activa el modo JSON forzado del modelo.
// timeout se aplica sobre ctx (el menor de los dos gana).
func (p *OllamaProvider) chat(ctx context.Context, system, user, format string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	body := ollamaChatReq{
		Model:  p.model,
		Stream: false,
		Messages: []oaiMsg{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Format: format,
	}

	data, err := doPost(ctx, p.baseURL+"/api/chat", nil, body)
	if err != nil {
		return "", fmt.Errorf("ollama http: %w", err)
	}

	var r ollamaChatResp
	if err := json.Unmarshal(data, &r); err != nil {
		return "", fmt.Errorf("ollama parse: %w | %.200s", err, data)
	}
	if r.Error != "" {
		return "", fmt.Errorf("ollama error: %s", r.Error)
	}
	return r.Message.Content, nil
}

// ── Helpers de verificación ───────────────────────────────────────────────────

type ollamaTagsResp struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
}

// ollamaListModels retorna los nombres de los modelos instalados en Ollama.
func ollamaListModels(baseURL string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/api/tags", nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tags ollamaTagsResp
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, fmt.Errorf("parse /api/tags: %w", err)
	}

	names := make([]string, 0, len(tags.Models))
	for _, m := range tags.Models {
		names = append(names, m.Name)
	}
	return names, nil
}

// ollamaModelExists comprueba si el modelo pedido (sin tag) está en la lista.
// "deepseek-r1" coincide con "deepseek-r1:latest" o "deepseek-r1:7b".
func ollamaModelExists(available []string, want string) bool {
	for _, name := range available {
		if name == want || strings.HasPrefix(name, want+":") {
			return true
		}
	}
	return false
}
