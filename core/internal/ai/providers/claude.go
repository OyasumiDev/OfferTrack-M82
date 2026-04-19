// core/internal/ai/providers/claude.go
package providers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/QUERTY/OfferTrack-M82/internal/domain"
	"github.com/QUERTY/OfferTrack-M82/internal/prompts"
)

const claudeBase = "https://api.anthropic.com/v1"

// ClaudeProvider implementa AIProvider usando la API de Anthropic Messages.
type ClaudeProvider struct{ apiKey, model string }

func NewClaudeProvider(apiKey, model string) (*ClaudeProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY no configurada")
	}
	if model == "" {
		model = "claude-3-5-haiku-20241022"
	}
	return &ClaudeProvider{apiKey: apiKey, model: model}, nil
}

func (p *ClaudeProvider) Name() string { return "claude" }

func (p *ClaudeProvider) complete(ctx context.Context, system, user string) (string, error) {
	body := map[string]any{
		"model":      p.model,
		"max_tokens": 4096,
		"system":     system,
		"messages":   []map[string]string{{"role": "user", "content": user}},
	}
	data, err := doPost(ctx, claudeBase+"/messages", map[string]string{
		"x-api-key":         p.apiKey,
		"anthropic-version": "2023-06-01",
	}, body)
	if err != nil {
		return "", fmt.Errorf("claude: %w", err)
	}

	var r struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Error *struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		} `json:"error,omitempty"`
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return "", fmt.Errorf("claude: parse → %w | %.200s", err, data)
	}
	if r.Error != nil {
		return "", fmt.Errorf("claude API (%s): %s", r.Error.Type, r.Error.Message)
	}
	if len(r.Content) == 0 {
		return "", fmt.Errorf("claude: respuesta vacía | %.200s", data)
	}
	return r.Content[0].Text, nil
}

func (p *ClaudeProvider) Analyze(ctx context.Context, req domain.AnalysisRequest) (*domain.AnalysisResult, error) {
	sys, user := prompts.BuildAnalysisMessages(req)
	raw, err := p.complete(ctx, sys, user)
	if err != nil {
		return nil, err
	}
	return prompts.ParseAnalysisResult(raw)
}

func (p *ClaudeProvider) AdaptCV(ctx context.Context, req domain.AdaptRequest) (*domain.AdaptResult, error) {
	sys, user := prompts.BuildAdaptMessages(req)
	raw, err := p.complete(ctx, sys, user)
	if err != nil {
		return nil, err
	}
	return prompts.ParseAdaptResult(raw)
}

func (p *ClaudeProvider) Summarize(ctx context.Context, text string) (string, error) {
	sys, user := prompts.BuildSummarizeMessages(text)
	return p.complete(ctx, sys, user)
}
