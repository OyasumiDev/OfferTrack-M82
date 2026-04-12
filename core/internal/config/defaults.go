package config

import "github.com/spf13/viper"

func SetDefaults() {
	viper.SetDefault("app.name", "OfferTrack M82")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.env", "development")
	viper.SetDefault("ai.provider", "gemini")
	viper.SetDefault("ai.model", "gemini-2.0-flash")
	viper.SetDefault("ai.timeout", 30)
	viper.SetDefault("qdrant.host", "localhost")
	viper.SetDefault("qdrant.port", 6334)
	viper.SetDefault("scraper.base_url", "http://localhost:3001")
	viper.SetDefault("scraper.timeout", 60)
}
