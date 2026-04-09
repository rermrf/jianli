package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"jianli/internal/middleware"
	"jianli/internal/store"
)

func setupSettingsTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := store.Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() returned error: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	settingsHandler := NewSettingsHandler(store.NewSiteSettingsStore(db))
	router := gin.New()
	protected := router.Group("/api")
	protected.Use(middleware.Auth("resume-key"))
	protected.GET("/settings", settingsHandler.Get)
	protected.PUT("/settings", settingsHandler.Update)

	return router
}

func TestSettingsHandlerGetRequiresAuth(t *testing.T) {
	router := setupSettingsTestRouter(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/settings", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestSettingsHandlerGetReturnsCurrentSettings(t *testing.T) {
	router := setupSettingsTestRouter(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	request.Header.Set("X-Auth-Key", "resume-key")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var payload struct {
		Data struct {
			AllowPDFExport bool `json:"allowPdfExport"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}

	if !payload.Data.AllowPDFExport {
		t.Fatal("expected allowPdfExport to default to true")
	}
}

func TestSettingsHandlerUpdatePersistsToggle(t *testing.T) {
	router := setupSettingsTestRouter(t)
	body := bytes.NewBufferString(`{"allowPdfExport":false}`)
	request := httptest.NewRequest(http.MethodPut, "/api/settings", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Key", "resume-key")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	getRecorder := httptest.NewRecorder()
	getRequest := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	getRequest.Header.Set("X-Auth-Key", "resume-key")
	router.ServeHTTP(getRecorder, getRequest)

	var payload struct {
		Data struct {
			AllowPDFExport bool `json:"allowPdfExport"`
		} `json:"data"`
	}
	if err := json.Unmarshal(getRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}

	if payload.Data.AllowPDFExport {
		t.Fatal("expected allowPdfExport to persist false")
	}
}
