package pdf

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildResumeViewIncludesAvatarEducationAwardsAndProjectURL(t *testing.T) {
	uploadRoot := filepath.Join(t.TempDir(), "uploads")
	avatarDir := filepath.Join(uploadRoot, "avatars")
	if err := os.MkdirAll(avatarDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	avatarPath := filepath.Join(avatarDir, "avatar.png")
	pngBytes := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xDD, 0x8D,
		0x18, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}
	if err := os.WriteFile(avatarPath, pngBytes, 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	resume := json.RawMessage(`{
		"profile": {
			"name": "Test User",
			"title": "Backend Engineer",
			"age": 25,
			"gender": "男",
			"education": "本科",
			"experience": "1年",
			"location": "Hangzhou",
			"hometown": "Ganzhou",
			"phone": "123456789",
			"email": "test@example.com",
			"avatarUrl": "/uploads/avatars/avatar.png"
		},
		"skills": ["Go", "MySQL"],
		"workExperience": [],
		"projects": [{
			"name": "AI Gateway",
			"startDate": "2025.12",
			"endDate": "2026.01",
			"description": ["Unified model gateway"],
			"url": "https://github.com/example/ai-gateway"
		}],
		"education": [{
			"school": "JXUFE",
			"major": "Computer Science",
			"degree": "Bachelor",
			"startDate": "2023.09",
			"endDate": "2025.07"
		}],
		"awards": [{
			"date": "2022.09",
			"title": "Scholarship"
		}]
	}`)

	view, err := buildResumeView(resume, uploadRoot)
	if err != nil {
		t.Fatalf("buildResumeView() returned error: %v", err)
	}

	html, err := renderResumeHTML(view)
	if err != nil {
		t.Fatalf("renderResumeHTML() returned error: %v", err)
	}

	if !strings.Contains(string(view.Profile.AvatarDataURL), "data:image/") {
		t.Fatalf("expected avatar data url in view, got %q", view.Profile.AvatarDataURL)
	}

	if len(view.Projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(view.Projects))
	}

	if view.Projects[0].URL != "https://github.com/example/ai-gateway" {
		t.Fatalf("expected project url to round-trip, got %q", view.Projects[0].URL)
	}

	for _, expected := range []string{
		"data:image/",
		"Test User",
		"Backend Engineer",
		"JXUFE",
		"Scholarship",
		"AI Gateway",
	} {
		if !strings.Contains(html, expected) {
			t.Fatalf("expected rendered html to include %q", expected)
		}
	}

	for _, unexpected := range []string{"print preview", "personal info"} {
		if strings.Contains(strings.ToLower(html), unexpected) {
			t.Fatalf("expected rendered html to avoid %q", unexpected)
		}
	}
}

func TestResolveBrowserPathPrefersConfiguredPath(t *testing.T) {
	configuredPath := filepath.Join(t.TempDir(), "chrome.exe")
	if err := os.WriteFile(configuredPath, []byte("stub"), 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	path, err := resolveBrowserPath(configuredPath, []string{"missing.exe"})
	if err != nil {
		t.Fatalf("resolveBrowserPath() returned error: %v", err)
	}

	if path != configuredPath {
		t.Fatalf("expected configured path %q, got %q", configuredPath, path)
	}
}

func TestRenderResumeHTMLUsesChineseCapableFontStack(t *testing.T) {
	html, err := renderResumeHTML(resumeView{})
	if err != nil {
		t.Fatalf("renderResumeHTML() returned error: %v", err)
	}

	for _, expected := range []string{
		"Noto Sans CJK SC",
		"WenQuanYi Micro Hei",
		"Microsoft YaHei",
	} {
		if !strings.Contains(html, expected) {
			t.Fatalf("expected rendered html to include font %q", expected)
		}
	}
}