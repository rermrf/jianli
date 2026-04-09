package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"jianli/internal/config"
	"jianli/internal/store"
)

func TestNewRouterRegistersCoreRoutes(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() returned error: %v", err)
	}
	defer sqlDB.Close()

	router := newRouter(config.Config{AuthKey: "resume-key"}, db)

	for _, testCase := range []struct {
		method string
		path   string
		codes  []int
	}{
		{method: http.MethodGet, path: "/api/resume", codes: []int{http.StatusOK}},
		{method: http.MethodGet, path: "/api/resume/pdf", codes: []int{http.StatusForbidden, http.StatusInternalServerError, http.StatusOK}},
		{method: http.MethodPost, path: "/api/auth/verify", codes: []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusOK}},
		{method: http.MethodGet, path: "/api/settings", codes: []int{http.StatusUnauthorized, http.StatusOK}},
		{method: http.MethodPut, path: "/api/settings", codes: []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusOK}},
		{method: http.MethodPost, path: "/api/upload/avatar", codes: []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusOK}},
		{method: http.MethodGet, path: "/api/visitors", codes: []int{http.StatusUnauthorized, http.StatusOK}},
		{method: http.MethodGet, path: "/api/visitors/stats", codes: []int{http.StatusUnauthorized, http.StatusOK}},
		{method: http.MethodPost, path: "/api/visitors/1", codes: []int{http.StatusBadRequest, http.StatusOK}},
		{method: http.MethodPost, path: "/api/resume/drafts", codes: []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusOK}},
		{method: http.MethodGet, path: "/api/resume/drafts", codes: []int{http.StatusUnauthorized, http.StatusOK}},
		{method: http.MethodGet, path: "/api/resume/drafts/1", codes: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusOK}},
		{method: http.MethodPut, path: "/api/resume/drafts/1/publish", codes: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusOK}},
		{method: http.MethodDelete, path: "/api/resume/drafts/1", codes: []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusOK}},
		{method: http.MethodGet, path: "/uploads/avatars/missing.png", codes: []int{http.StatusNotFound, http.StatusOK}},
	} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(testCase.method, testCase.path, nil)
		router.ServeHTTP(recorder, request)

		matched := false
		for _, code := range testCase.codes {
			if recorder.Code == code {
				matched = true
				break
			}
		}

		if !matched {
			t.Fatalf("expected route %s %s to be registered, got status %d", testCase.method, testCase.path, recorder.Code)
		}
	}
}

func TestNewRouterServesFrontendIndexForSPAPaths(t *testing.T) {
	tempRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tempRoot, "web", "dist", "assets"), 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempRoot, "web", "dist", "index.html"), []byte("<!doctype html><html><body>jianli</body></html>"), 0o644); err != nil {
		t.Fatalf("WriteFile(index.html) returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempRoot, "web", "dist", "assets", "app.js"), []byte("console.log('jianli')"), 0o644); err != nil {
		t.Fatalf("WriteFile(app.js) returned error: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() returned error: %v", err)
	}
	if err := os.Chdir(tempRoot); err != nil {
		t.Fatalf("Chdir(tempRoot) returned error: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	}()

	db, err := store.Open(filepath.Join(tempRoot, "data", "resume.db"))
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() returned error: %v", err)
	}
	defer sqlDB.Close()

	router := newRouter(config.Config{AuthKey: "resume-key"}, db)

	for _, testCase := range []struct {
		path         string
		contentType  string
		responseBody string
	}{
		{path: "/", contentType: "text/html; charset=utf-8", responseBody: "jianli"},
		{path: "/drafts/123", contentType: "text/html; charset=utf-8", responseBody: "jianli"},
		{path: "/assets/app.js", contentType: "text/javascript; charset=utf-8", responseBody: "console.log('jianli')"},
	} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, testCase.path, nil)
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected GET %s to return 200, got %d", testCase.path, recorder.Code)
		}
		if recorder.Header().Get("Content-Type") != testCase.contentType {
			t.Fatalf("expected GET %s content-type %q, got %q", testCase.path, testCase.contentType, recorder.Header().Get("Content-Type"))
		}
		if !strings.Contains(recorder.Body.String(), testCase.responseBody) {
			t.Fatalf("expected GET %s body to contain %q, got %q", testCase.path, testCase.responseBody, recorder.Body.String())
		}
	}
}
