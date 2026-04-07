package pdf

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildResumeViewIncludesAvatarEducationAndAwards(t *testing.T) {
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
			"name": "温庆京",
			"title": "Golang 后端工程师",
			"location": "浙江杭州",
			"phone": "17620096266",
			"email": "3219431643@qq.com",
			"avatarUrl": "/uploads/avatars/avatar.png"
		},
		"skills": ["Go", "MySQL"],
		"workExperience": [],
		"projects": [],
		"education": [{
			"school": "江西财经大学现代经济管理学院",
			"major": "计算机科学与技术",
			"degree": "本科"
		}],
		"awards": [{
			"date": "2022.09",
			"title": "国家励志奖学金"
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

	if !strings.Contains(html, "data:image/") {
		t.Fatalf("expected rendered html to inline avatar data url, got %q", html)
	}

	if !strings.Contains(html, "教育经历") {
		t.Fatal("expected rendered html to include education section")
	}

	if !strings.Contains(html, "江西财经大学现代经济管理学院") {
		t.Fatal("expected rendered html to include education content")
	}

	if !strings.Contains(html, "荣誉奖项") {
		t.Fatal("expected rendered html to include awards section")
	}

	if !strings.Contains(html, "国家励志奖学金") {
		t.Fatal("expected rendered html to include awards content")
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
