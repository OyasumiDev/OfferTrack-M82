// core/internal/ai/providers/openrouter.go
package providers

import (
	"context"
	"fmt"

	"github.com/QUERTY/OfferTrack-M82/internal/domain"
	"github.com/QUERTY/OfferTrack-M82/internal/prompts"
)

const openrouterBase = "https://openrouter.ai/api/v1"

// OpenrouterProvider implementa AIProvider usando la API de OpenRouter (compatible con OpenAI).
type OpenrouterProvider struct{ apiKey, model string }

func NewOpenrouterProvider(apiKey, model string) (*OpenrouterProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY no configurada")
	}
	if model == "" {
		model = "google/gemini-2.0-flash-001"
	}
	return &OpenrouterProvider{apiKey: apiKey, model: model}, nil
}

func (p *OpenrouterProvider) Name() string { return "openrouter" }

func (p *OpenrouterProvider) Analyze(ctx context.Context, req domain.AnalysisRequest) (*domain.AnalysisResult, error) {
	sys, user := prompts.BuildAnalysisMessages(req)
	raw, err := doOpenAIChat(ctx, openrouterBase, p.apiKey, p.model, sys, user)
	if err != nil {
		return nil, fmt.Errorf("openrouter: %w", err)
	}
	return prompts.ParseAnalysisResult(raw)
}

func (p *OpenrouterProvider) AdaptCV(ctx context.Context, req domain.AdaptRequest) (*domain.AdaptResult, error) {
	sys, user := prompts.BuildAdaptMessages(req)
	raw, err := doOpenAIChat(ctx, openrouterBase, p.apiKey, p.model, sys, user)
	if err != nil {
		return nil, fmt.Errorf("openrouter: %w", err)
	}
	return prompts.ParseAdaptResult(raw)
}

func (p *OpenrouterProvider) Summarize(ctx context.Context, text string) (string, error) {
	sys, user := prompts.BuildSummarizeMessages(text)
	raw, err := doOpenAIChat(ctx, openrouterBase, p.apiKey, p.model, sys, user)
	if err != nil {
		return "", fmt.Errorf("openrouter: %w", err)
	}
	return raw, nil
}
