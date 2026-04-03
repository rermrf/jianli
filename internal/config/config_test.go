package config

import "testing"

func TestLoadUsesDefaultsAndEnvOverrides(t *testing.T) {
	t.Setenv("AUTH_KEY", "resume-key")
	t.Setenv("PORT", "9090")
	t.Setenv("DB_PATH", "./tmp/resume.db")
	t.Setenv("FRONTEND_ORIGIN", "http://localhost:5173")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Port != "9090" {
		t.Fatalf("expected port 9090, got %q", cfg.Port)
	}

	if cfg.AuthKey != "resume-key" {
		t.Fatalf("expected auth key to be loaded from env, got %q", cfg.AuthKey)
	}

	if cfg.DBPath != "./tmp/resume.db" {
		t.Fatalf("expected db path override, got %q", cfg.DBPath)
	}

	if cfg.FrontendOrigin != "http://localhost:5173" {
		t.Fatalf("expected frontend origin override, got %q", cfg.FrontendOrigin)
	}
}

func TestLoadRequiresAuthKey(t *testing.T) {
	t.Setenv("AUTH_KEY", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when AUTH_KEY is missing")
	}
}
