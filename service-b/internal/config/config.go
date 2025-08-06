package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port 			string
	WeatherAPIKey 	string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port: 			getEnv("PORT", "8080"),
		WeatherAPIKey:	getEnv("WEATHER_API_KEY", ""),
	}
	if cfg.WeatherAPIKey == "" {
		return nil, fmt.Errorf("WEATHER_API_KEY is required")
	}
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}