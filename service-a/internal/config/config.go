package config

import (
	"os"
)

type Config struct {
	Port 		string
	ServiceBURL string
}

func LoadConfig() (*Config, error) {
	cfg := &Config {
		Port: getEnv("PORT", "8081"),
		ServiceBURL: getEnv("SERVICE_B_URL", "http://service-b:8080"),
	}
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
