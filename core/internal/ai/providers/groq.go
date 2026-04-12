package providers

import (
	"context"
	"fmt"

	"github.com/yourusername/offertrack-m82/core/internal/ai"
)

type GroqProvider struct{ apiKey, model string }

func NewGroqProvider(apiKey, model string) (*GroqProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY no configurada")
	}
	return &GroqProvider{apiKey: apiKey, model: model}, nil
}

func (p *GroqProvider) Name() string { return "groq" }

func (p *GroqProvider) Analyze(ctx context.Context, req ai.AnalysisRequest) (*ai.AnalysisResult, error) {
	return nil, nil // TODO: implementar
}

func (p *GroqProvider) AdaptCV(ctx context.Context, req ai.AdaptRequest) (*ai.AdaptResult, error) {
	return nil, nil // TODO: implementar
}

func (p *GroqProvider) Summarize(ctx context.Context, text string) (string, error) {
	return "", nil // TODO: implementar
}
