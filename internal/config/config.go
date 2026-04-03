package config

import (
	"errors"
	"os"
)

type Config struct {
	AuthKey        string
	DBPath         string
	FrontendOrigin string
	Port           string
}

func Load() (Config, error) {
	cfg := Config{
		AuthKey:        os.Getenv("AUTH_KEY"),
		DBPath:         getenvOrDefault("DB_PATH", "./data/resume.db"),
		FrontendOrigin: getenvOrDefault("FRONTEND_ORIGIN", "http://localhost:5173"),
		Port:           getenvOrDefault("PORT", "8080"),
	}

	if cfg.AuthKey == "" {
		return Config{}, errors.New("AUTH_KEY is required")
	}

	return cfg, nil
}

func getenvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
