package ai

import (
	"fmt"

	"github.com/yourusername/offertrack-m82/core/internal/ai/providers"
	"github.com/yourusername/offertrack-m82/core/internal/config"
)

// NewProvider crea el proveedor correcto segun la configuracion.
// Agregar uno nuevo = archivo en providers/ + un case aqui.
func NewProvider(cfg *config.Config) (AIProvider, error) {
	switch cfg.AI.Provider {
	case "claude":
		return providers.NewClaudeProvider(cfg.AI.APIKey, cfg.AI.Model)
	case "gemini":
		return providers.NewGeminiProvider(cfg.AI.APIKey, cfg.AI.Model)
	case "groq":
		return providers.NewGroqProvider(cfg.AI.APIKey, cfg.AI.Model)
	case "openrouter":
		return providers.NewOpenRouterProvider(cfg.AI.APIKey, cfg.AI.Model)
	default:
		return nil, fmt.Errorf("proveedor IA no soportado: %s", cfg.AI.Provider)
	}
}
