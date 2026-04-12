package config

import "github.com/spf13/viper"

type Config struct {
	App     AppConfig
	AI      AIConfig
	Qdrant  QdrantConfig
	Scraper ScraperConfig
}

type AppConfig   struct{ Name, Version, Env string }
type AIConfig    struct{ Provider, Model, APIKey string; Timeout int }
type QdrantConfig struct{ Host string; Port int }
type ScraperConfig struct{ BaseURL string; Timeout int }

func Load() (*Config, error) {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	cfg := &Config{}
	return cfg, viper.Unmarshal(cfg)
}
