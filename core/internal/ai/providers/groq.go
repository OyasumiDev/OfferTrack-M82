// core/internal/ai/providers/groq.go
package providers

import (
	"context"
	"fmt"

	"github.com/QUERTY/OfferTrack-M82/internal/domain"
	"github.com/QUERTY/OfferTrack-M82/internal/prompts"
)

const groqBase = "https://api.groq.com/openai/v1"

// GroqProvider implementa AIProvider usando la API de Groq (compatible con OpenAI).
type GroqProvider struct{ apiKey, model string }

func NewGroqProvider(apiKey, model string) (*GroqProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY no configurada")
	}
	if model == "" {
		model = "llama-3.1-70b-versatile"
	}
	return &GroqProvider{apiKey: apiKey, model: model}, nil
}

func (p *GroqProvider) Name() string { return "groq" }

func (p *GroqProvider) Analyze(ctx context.Context, req domain.AnalysisRequest) (*domain.AnalysisResult, error) {
	sys, user := prompts.BuildAnalysisMessages(req)
	raw, err := doOpenAIChat(ctx, groqBase, p.apiKey, p.model, sys, user)
	if err != nil {
		return nil, fmt.Errorf("groq: %w", err)
	}
	return prompts.ParseAnalysisResult(raw)
}

func (p *GroqProvider) AdaptCV(ctx context.Context, req domain.AdaptRequest) (*domain.AdaptResult, error) {
	sys, user := prompts.BuildAdaptMessages(req)
	raw, err := doOpenAIChat(ctx, groqBase, p.apiKey, p.model, sys, user)
	if err != nil {
		return nil, fmt.Errorf("groq: %w", err)
	}
	return prompts.ParseAdaptResult(raw)
}

func (p *GroqProvider) Summarize(ctx context.Context, text string) (string, error) {
	sys, user := prompts.BuildSummarizeMessages(text)
	raw, err := doOpenAIChat(ctx, groqBase, p.apiKey, p.model, sys, user)
	if err != nil {
		return "", fmt.Errorf("groq: %w", err)
	}
	return raw, nil
}
