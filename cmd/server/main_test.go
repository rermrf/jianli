package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"jianli/internal/config"
	"jianli/internal/model"
	"jianli/internal/store"
)

func TestNewRouterRecordsFullClientIPThroughDockerProxy(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	router := newRouter(config.Config{
		AuthKey: "resume-key", TrustedProxies: []string{"127.0.0.1", "::1", "172.16.0.0/12"},
	}, db)
	for _, tc := range []struct{ name, remote, forwarded, wantIP string }{
		{"Docker proxy chain", "172.20.0.3:54321", "112.17.45.201, 172.20.0.4", "112.17.45.201"},
		{"forwarded IPv6", "172.20.0.3:54321", "240e:3b7:3272:d8d0:db09:c067:8d59:539e", "240e:3b7:3272:d8d0:db09:c067:8d59:539e"},
		{"untrusted public sender", "8.8.8.8:54321", "112.17.45.202", "8.8.8.8"},
		{"direct IPv6", "[2606:4700:4700::1111]:54321", "", "2606:4700:4700::1111"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/visitors", strings.NewReader(`{"duration":0}`))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-Forwarded-For", tc.forwarded)
			request.RemoteAddr = tc.remote
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusOK {
				t.Fatalf("create status=%d body=%s", recorder.Code, recorder.Body.String())
			}
			var payload struct {
				Data struct{ ID int64 } `json:"data"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
				t.Fatal(err)
			}
			var record model.VisitorRecord
			if err := db.First(&record, payload.Data.ID).Error; err != nil {
				t.Fatal(err)
			}
			if record.IP != tc.wantIP {
				t.Fatalf("stored IP=%q, want full client IP %q", record.IP, tc.wantIP)
			}
		})
	}
}

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
