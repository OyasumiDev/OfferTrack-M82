package providers

import (
	"context"
	"fmt"

	"github.com/yourusername/offertrack-m82/core/internal/ai"
)

type GeminiProvider struct{ apiKey, model string }

func NewGeminiProvider(apiKey, model string) (*GeminiProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY no configurada")
	}
	return &GeminiProvider{apiKey: apiKey, model: model}, nil
}

func (p *GeminiProvider) Name() string { return "gemini" }

func (p *GeminiProvider) Analyze(ctx context.Context, req ai.AnalysisRequest) (*ai.AnalysisResult, error) {
	return nil, nil // TODO: implementar
}

func (p *GeminiProvider) AdaptCV(ctx context.Context, req ai.AdaptRequest) (*ai.AdaptResult, error) {
	return nil, nil // TODO: implementar
}

func (p *GeminiProvider) Summarize(ctx context.Context, text string) (string, error) {
	return "", nil // TODO: implementar
}
