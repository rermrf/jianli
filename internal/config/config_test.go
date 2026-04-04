package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUsesConfigFileValues(t *testing.T) {
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

	_, err = Load()
	if err == nil {
		t.Fatal("expected error when config.json is missing")
	}
}

func TestLoadRequiresAuthKeyInConfigFile(t *testing.T) {
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

	if err := os.WriteFile("config.json", []byte(`{"port":"8080"}`), 0o644); err != nil {
		t.Fatalf("WriteFile(config.json) returned error: %v", err)
	}

	_, err = Load()
	if err == nil {
		t.Fatal("expected error when authKey is missing")
	}
}

func TestLoadSupportsUTF8BOMConfigFile(t *testing.T) {
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
