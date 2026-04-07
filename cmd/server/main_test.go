package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
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
		{method: http.MethodGet, path: "/api/resume/pdf", codes: []int{http.StatusInternalServerError, http.StatusOK}},
		{method: http.MethodPost, path: "/api/auth/verify", codes: []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusOK}},
		{method: http.MethodPost, path: "/api/upload/avatar", codes: []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusOK}},
		{method: http.MethodGet, path: "/api/visitors", codes: []int{http.StatusUnauthorized, http.StatusOK}},
		{method: http.MethodGet, path: "/api/visitors/stats", codes: []int{http.StatusUnauthorized, http.StatusOK}},
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
