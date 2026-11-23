package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
)

type Config struct {
	AuthKey         string   `json:"authKey"`
	BrowserPath     string   `json:"browserPath"`
	DBPath          string   `json:"dbPath"`
	FrontendOrigin  string   `json:"frontendOrigin"`
	IP2RegionPath   string   `json:"ip2RegionPath"`
	IP2RegionV6Path string   `json:"ip2RegionV6Path"`
	TrustedProxies  []string `json:"trustedProxies"`
	Port            string   `json:"port"`
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
	if cfg.IP2RegionV6Path == "" {
		cfg.IP2RegionV6Path = "./data/ip2region_v6.xdb"
	}
	if cfg.TrustedProxies == nil {
		// Docker bridge networks normally allocate addresses from 172.16/12.
		cfg.TrustedProxies = []string{"127.0.0.1", "::1", "172.16.0.0/12"}
	}
	for _, proxy := range cfg.TrustedProxies {
		if net.ParseIP(proxy) == nil {
			if _, _, err := net.ParseCIDR(proxy); err != nil {
				return Config{}, fmt.Errorf("invalid trustedProxies entry %q: %w", proxy, err)
			}
		}
	}

	return cfg, nil
}
