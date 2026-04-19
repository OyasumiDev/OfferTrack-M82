// core/internal/config/config.go
package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig
	AI      AIConfig
	Qdrant  QdrantConfig
	Scraper ScraperConfig
}

type AppConfig struct {
	Name    string
	Version string
	Env     string
}

type AIConfig struct {
	Provider      string
	Model         string
	APIKey        string
	Timeout       int
	OllamaBaseURL string // solo se usa cuando Provider = "ollama"
}

type QdrantConfig struct {
	Host        string
	Port        int
	Collections map[string]string // keys: jobs, profile, cvs, memory
}

type ScraperConfig struct {
	BaseURL string
	Timeout int
}

// Load carga la configuración combinando app.yaml y variables de entorno.
// Las variables de entorno tienen prioridad sobre el yaml.
func Load() (*Config, error) {
	SetDefaults()

	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AutomaticEnv()

	// Ignorar error si no existe app.yaml — los defaults y el .env son suficientes
	_ = viper.ReadInConfig()

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Las colecciones siempre se leen desde variables de entorno
	cfg.Qdrant.Collections = map[string]string{
		"jobs":    envOrDefault("QDRANT_COLLECTION_JOBS", "jobs"),
		"profile": envOrDefault("QDRANT_COLLECTION_PROFILE", "profile"),
		"cvs":     envOrDefault("QDRANT_COLLECTION_CVS", "cvs"),
		"memory":  envOrDefault("QDRANT_COLLECTION_MEMORY", "claude_memory"),
	}

	// El puerto de Qdrant también puede venir del .env
	if portStr := os.Getenv("QDRANT_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Qdrant.Port = port
		}
	}
	if host := os.Getenv("QDRANT_HOST"); host != "" {
		cfg.Qdrant.Host = host
	}

	// Leer AI_PROVIDER y AI_MODEL directamente del entorno — viper no mapea ai.provider → AI_PROVIDER
	// sin SetEnvKeyReplacer, así que lo hacemos explícitamente.
	if v := os.Getenv("AI_PROVIDER"); v != "" {
		cfg.AI.Provider = v
	}
	if v := os.Getenv("AI_MODEL"); v != "" {
		cfg.AI.Model = v
	}

	// API key según el proveedor activo
	cfg.AI.APIKey = resolveAPIKey(cfg.AI.Provider)

	// Base URL del scraper
	if u := os.Getenv("SCRAPER_BASE_URL"); u != "" {
		cfg.Scraper.BaseURL = u
	}

	// Base URL de Ollama (solo relevante si AI_PROVIDER=ollama)
	cfg.AI.OllamaBaseURL = envOrDefault("OLLAMA_BASE_URL", "http://localhost:11434")

	return cfg, nil
}

func resolveAPIKey(provider string) string {
	switch provider {
	case "claude":
		return os.Getenv("ANTHROPIC_API_KEY")
	case "gemini":
		return os.Getenv("GEMINI_API_KEY")
	case "groq":
		return os.Getenv("GROQ_API_KEY")
	case "openrouter":
		return os.Getenv("OPENROUTER_API_KEY")
	default:
		return ""
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
