package providers

import (
	"context"
	"fmt"

	"github.com/yourusername/offertrack-m82/core/internal/ai"
)

type OpenrouterProvider struct{ apiKey, model string }

func NewOpenrouterProvider(apiKey, model string) (*OpenrouterProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY no configurada")
	}
	return &OpenrouterProvider{apiKey: apiKey, model: model}, nil
}

func (p *OpenrouterProvider) Name() string { return "openrouter" }

func (p *OpenrouterProvider) Analyze(ctx context.Context, req ai.AnalysisRequest) (*ai.AnalysisResult, error) {
	return nil, nil // TODO: implementar
}

func (p *OpenrouterProvider) AdaptCV(ctx context.Context, req ai.AdaptRequest) (*ai.AdaptResult, error) {
	return nil, nil // TODO: implementar
}

func (p *OpenrouterProvider) Summarize(ctx context.Context, text string) (string, error) {
	return "", nil // TODO: implementar
}
