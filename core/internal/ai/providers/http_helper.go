// core/internal/ai/providers/http_helper.go
package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 120 * time.Second}

// openAIReq es el cuerpo de las APIs compatibles con OpenAI /chat/completions.
type openAIReq struct {
	Model     string     `json:"model"`
	Messages  []oaiMsg   `json:"messages"`
	MaxTokens int        `json:"max_tokens,omitempty"`
}

type oaiMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// doOpenAIChat llama a un endpoint compatible con OpenAI /chat/completions.
func doOpenAIChat(ctx context.Context, baseURL, apiKey, model, system, user string) (string, error) {
	payload := openAIReq{
		Model:     model,
		MaxTokens: 4096,
		Messages: []oaiMsg{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
	}
	data, err := doPost(ctx, baseURL+"/chat/completions", map[string]string{
		"Authorization": "Bearer " + apiKey,
	}, payload)
	if err != nil {
		return "", err
	}

	var r openAIResp
	if err := json.Unmarshal(data, &r); err != nil {
		return "", fmt.Errorf("openai: parse → %w | %.200s", err, data)
	}
	if r.Error != nil {
		return "", fmt.Errorf("openai API: %s", r.Error.Message)
	}
	if len(r.Choices) == 0 {
		return "", fmt.Errorf("openai: sin respuesta | %.200s", data)
	}
	return r.Choices[0].Message.Content, nil
}

// doPost ejecuta un POST JSON y retorna el body raw.
func doPost(ctx context.Context, url string, headers map[string]string, body any) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("doPost: marshal → %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("doPost: request → %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doPost: http → %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("doPost: read body → %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("doPost: status %d | %.300s", resp.StatusCode, data)
	}
	return data, nil
}
