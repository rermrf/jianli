package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"jianli/internal/middleware"
	"jianli/internal/store"
)

func setupResumeTestRouter(t *testing.T) *gin.Engine {
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

	resumeHandler := NewResumeHandler(store.NewResumeStore(db))
	router := gin.New()
	router.GET("/api/resume", resumeHandler.Get)
	protected := router.Group("/api")
	protected.Use(middleware.Auth("resume-key"))
	protected.PUT("/resume", resumeHandler.Update)

	return router
}

func TestResumeHandlerGetReturnsSeededResume(t *testing.T) {
	router := setupResumeTestRouter(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/resume", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var payload struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}

	profile := payload.Data["profile"].(map[string]any)
	if profile["name"] != "温庆京" {
		t.Fatalf("expected seeded profile name 温庆京, got %#v", profile["name"])
	}
}

func TestResumeHandlerUpdateRequiresAuthAndSavesData(t *testing.T) {
	router := setupResumeTestRouter(t)
	updated := `{"profile":{"name":"测试后端"}}`

	request := httptest.NewRequest(http.MethodPut, "/api/resume", strings.NewReader(updated))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Key", "resume-key")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	getRecorder := httptest.NewRecorder()
	getRequest := httptest.NewRequest(http.MethodGet, "/api/resume", nil)
	router.ServeHTTP(getRecorder, getRequest)

	var payload struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(getRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}

	profile := payload.Data["profile"].(map[string]any)
	if profile["name"] != "测试后端" {
		t.Fatalf("expected updated profile name 测试后端, got %#v", profile["name"])
	}
}
