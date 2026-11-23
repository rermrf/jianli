package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func withTempConfigDir(t *testing.T) string {
	t.Helper()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() returned error: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir(%q) returned error: %v", tempDir, err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(originalDir)
	})

	return tempDir
}

func TestLoadDefaultsSupportDockerProxiesAndBothIPFamilies(t *testing.T) {
	withTempConfigDir(t)
	if err := os.WriteFile("config.json", []byte(`{"authKey":"resume-key"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.IP2RegionPath != "./data/ip2region_v4.xdb" || cfg.IP2RegionV6Path != "./data/ip2region_v6.xdb" {
		t.Fatalf("unexpected database paths: %q %q", cfg.IP2RegionPath, cfg.IP2RegionV6Path)
	}
	if !reflect.DeepEqual(cfg.TrustedProxies, []string{"127.0.0.1", "::1", "172.16.0.0/12"}) {
		t.Fatalf("unexpected trusted proxies: %v", cfg.TrustedProxies)
	}
}

func TestLoadAllowsDisablingProxiesAndRejectsInvalidProxies(t *testing.T) {
	withTempConfigDir(t)
	for _, tc := range []struct {
		proxies string
		wantErr bool
	}{
		{`[]`, false},
		{`["192.168.20.0/24", "::1"]`, false},
		{`["*"]`, true},
	} {
		if err := os.WriteFile("config.json", []byte(`{"authKey":"resume-key","trustedProxies":`+tc.proxies+`}`), 0o644); err != nil {
			t.Fatal(err)
		}
		cfg, err := Load()
		if (err != nil) != tc.wantErr {
			t.Fatalf("proxies %s: err=%v", tc.proxies, err)
		}
		if tc.proxies == `[]` && len(cfg.TrustedProxies) != 0 {
			t.Fatal("explicit empty trustedProxies must disable forwarded headers")
		}
	}
}

func TestLoadUsesAuthKeyFromConfigFile(t *testing.T) {
	tempDir := withTempConfigDir(t)

	configPath := filepath.Join(tempDir, "config.json")
	configContent := `{
  "authKey": "resume-key",
  "port": "9090",
  "dbPath": "./tmp/resume.db",
  "frontendOrigin": "http://localhost:5173"
}`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) returned error: %v", configPath, err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Port != "9090" {
		t.Fatalf("expected port 9090, got %q", cfg.Port)
	}

	if cfg.AuthKey != "resume-key" {
		t.Fatalf("expected auth key from config file, got %q", cfg.AuthKey)
	}

	if cfg.DBPath != "./tmp/resume.db" {
		t.Fatalf("expected db path from config file, got %q", cfg.DBPath)
	}

	if cfg.FrontendOrigin != "http://localhost:5173" {
		t.Fatalf("expected frontend origin from config file, got %q", cfg.FrontendOrigin)
	}
}

func TestLoadRequiresConfigFile(t *testing.T) {
	withTempConfigDir(t)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when config.json is missing")
	}
}

func TestLoadRequiresAuthKeyInConfigFile(t *testing.T) {
	withTempConfigDir(t)

	if err := os.WriteFile("config.json", []byte(`{"port":"8080"}`), 0o644); err != nil {
		t.Fatalf("WriteFile(config.json) returned error: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when authKey is missing")
	}
}

func TestLoadSupportsUTF8BOMConfigFile(t *testing.T) {
	withTempConfigDir(t)

	content := append([]byte{0xEF, 0xBB, 0xBF}, []byte(`{"authKey":"resume-key","port":"8080"}`)...)
	if err := os.WriteFile("config.json", content, 0o644); err != nil {
		t.Fatalf("WriteFile(config.json) returned error: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error for BOM config: %v", err)
	}

	if cfg.AuthKey != "resume-key" {
		t.Fatalf("expected authKey from BOM config, got %q", cfg.AuthKey)
	}
}
