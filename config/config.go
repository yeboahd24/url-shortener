package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	PostgresDSN  string `mapstructure:"POSTGRES_DSN"`
	RedisAddr    string `mapstructure:"REDIS_ADDR"`
	RedisPass    string `mapstructure:"REDIS_PASS"`
	GeoAPIURL    string `mapstructure:"GEO_API_URL"`
	Port         string `mapstructure:"PORT"`
	APIKeyHeader string `mapstructure:"API_KEY_HEADER"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("POSTGRES_DSN", "postgres://user:pass@localhost:5432/urlshortener")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASS", "")
	viper.SetDefault("GEO_API_URL", "http://ip-api.com/json")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("API_KEY_HEADER", "X-API-Key")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
