package providers

import (
	"context"
	"fmt"

	"github.com/yourusername/offertrack-m82/core/internal/ai"
)

type ClaudeProvider struct{ apiKey, model string }

func NewClaudeProvider(apiKey, model string) (*ClaudeProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY no configurada")
	}
	return &ClaudeProvider{apiKey: apiKey, model: model}, nil
}

func (p *ClaudeProvider) Name() string { return "claude" }

func (p *ClaudeProvider) Analyze(ctx context.Context, req ai.AnalysisRequest) (*ai.AnalysisResult, error) {
	return nil, nil // TODO: implementar
}

func (p *ClaudeProvider) AdaptCV(ctx context.Context, req ai.AdaptRequest) (*ai.AdaptResult, error) {
	return nil, nil // TODO: implementar
}

func (p *ClaudeProvider) Summarize(ctx context.Context, text string) (string, error) {
	return "", nil // TODO: implementar
}
