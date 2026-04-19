// core/internal/ai/providers/gemini.go
package providers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/QUERTY/OfferTrack-M82/internal/domain"
	"github.com/QUERTY/OfferTrack-M82/internal/prompts"
)

// GeminiProvider implementa AIProvider usando la API REST de Google Gemini.
type GeminiProvider struct{ apiKey, model string }

func NewGeminiProvider(apiKey, model string) (*GeminiProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY no configurada")
	}
	if model == "" {
		model = "gemini-2.0-flash"
	}
	return &GeminiProvider{apiKey: apiKey, model: model}, nil
}

func (p *GeminiProvider) Name() string { return "gemini" }

func (p *GeminiProvider) complete(ctx context.Context, system, user string) (string, error) {
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		p.model, p.apiKey,
	)
	body := map[string]any{
		"systemInstruction": map[string]any{
			"parts": []map[string]string{{"text": system}},
		},
		"contents": []map[string]any{
			{
				"role":  "user",
				"parts": []map[string]string{{"text": user}},
			},
		},
		"generationConfig": map[string]any{
			"maxOutputTokens": 4096,
		},
	}
	data, err := doPost(ctx, url, nil, body)
	if err != nil {
		return "", fmt.Errorf("gemini: %w", err)
	}

	var r struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		Error *struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		} `json:"error,omitempty"`
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return "", fmt.Errorf("gemini: parse → %w | %.200s", err, data)
	}
	if r.Error != nil {
		return "", fmt.Errorf("gemini API (%d): %s", r.Error.Code, r.Error.Message)
	}
	if len(r.Candidates) == 0 || len(r.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini: respuesta vacía | %.200s", data)
	}
	return r.Candidates[0].Content.Parts[0].Text, nil
}

func (p *GeminiProvider) Analyze(ctx context.Context, req domain.AnalysisRequest) (*domain.AnalysisResult, error) {
	sys, user := prompts.BuildAnalysisMessages(req)
	raw, err := p.complete(ctx, sys, user)
	if err != nil {
		return nil, err
	}
	return prompts.ParseAnalysisResult(raw)
}

func (p *GeminiProvider) AdaptCV(ctx context.Context, req domain.AdaptRequest) (*domain.AdaptResult, error) {
	sys, user := prompts.BuildAdaptMessages(req)
	raw, err := p.complete(ctx, sys, user)
	if err != nil {
		return nil, err
	}
	return prompts.ParseAdaptResult(raw)
}

func (p *GeminiProvider) Summarize(ctx context.Context, text string) (string, error) {
	sys, user := prompts.BuildSummarizeMessages(text)
	return p.complete(ctx, sys, user)
}
