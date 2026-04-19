package ai

import (
	"fmt"

	"github.com/QUERTY/OfferTrack-M82/internal/ai/providers"
	"github.com/QUERTY/OfferTrack-M82/internal/config"
)

// NewProvider crea el proveedor correcto según la configuración.
// Agregar uno nuevo = archivo en providers/ + un case aquí.
func NewProvider(cfg *config.Config) (AIProvider, error) {
	switch cfg.AI.Provider {
	case "claude":
		return providers.NewClaudeProvider(cfg.AI.APIKey, cfg.AI.Model)
	case "gemini":
		return providers.NewGeminiProvider(cfg.AI.APIKey, cfg.AI.Model)
	case "groq":
		return providers.NewGroqProvider(cfg.AI.APIKey, cfg.AI.Model)
	case "openrouter":
		return providers.NewOpenrouterProvider(cfg.AI.APIKey, cfg.AI.Model)
	case "ollama":
		model := cfg.AI.Model
		if model == "" {
			model = "deepseek-r1"
		}
		return providers.NewOllamaProvider(cfg.AI.OllamaBaseURL, model)
	default:
		return nil, fmt.Errorf("proveedor IA no soportado: %q\n(valores válidos: claude, gemini, groq, openrouter, ollama)", cfg.AI.Provider)
	}
}
