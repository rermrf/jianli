package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	AuthKey        string `json:"authKey"`
	BrowserPath    string `json:"browserPath"`
	DBPath         string `json:"dbPath"`
	FrontendOrigin string `json:"frontendOrigin"`
	IP2RegionPath  string `json:"ip2RegionPath"`
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

	if cfg.AuthKey == "" {
		return Config{}, errors.New("authKey is required in config.json")
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
	if cfg.IP2RegionPath == "" {
		cfg.IP2RegionPath = "./data/ip2region_v4.xdb"
	}

	return cfg, nil
}
