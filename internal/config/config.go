package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	AuthKey        string `json:"-"`
	BrowserPath    string `json:"browserPath"`
	DBPath         string `json:"dbPath"`
	FrontendOrigin string `json:"frontendOrigin"`
	Port           string `json:"port"`
}

func Load() (Config, error) {
	content, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, fmt.Errorf("read config.json: %w", err)
	}

	content = bytes.TrimPrefix(content, []byte{0xEF, 0xBB, 0xBF})

	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode config.json: %w", err)
	}

	cfg.AuthKey = os.Getenv("AUTH_KEY")
	if cfg.AuthKey == "" {
		return Config{}, errors.New("AUTH_KEY is required")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DBPath == "" {
		cfg.DBPath = "./data/resume.db"
	}
	if cfg.FrontendOrigin == "" {
		cfg.FrontendOrigin = "http://localhost:5173"
	}

	return cfg, nil
}
