package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"jianli/internal/store"
)

func setupVisitorTestRouter(t *testing.T) *gin.Engine {
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

	visitorHandler := NewVisitorHandler(store.NewVisitorStore(db))
	router := gin.New()
	router.POST("/api/visitors", visitorHandler.Create)
	router.GET("/api/visitors", visitorHandler.List)
	router.GET("/api/visitors/stats", visitorHandler.Stats)
	router.PATCH("/api/visitors/:id", visitorHandler.UpdateDuration)

	return router
}

func TestVisitorHandlerCreateAndList(t *testing.T) {
	router := setupVisitorTestRouter(t)

	request := httptest.NewRequest(http.MethodPost, "/api/visitors", strings.NewReader(`{"ip":"112.17.45.201","city":"杭州","device":"Windows","browser":"Chrome","os":"Windows","visitTime":"2026-04-03T10:00:00Z","duration":0}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 from create, got %d", recorder.Code)
	}

	listRecorder := httptest.NewRecorder()
	listRequest := httptest.NewRequest(http.MethodGet, "/api/visitors?days=7&page=1&limit=20", nil)
	router.ServeHTTP(listRecorder, listRequest)

	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 from list, got %d", listRecorder.Code)
	}

	var payload struct {
		Code int `json:"code"`
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected list JSON, got error: %v", err)
	}

	if len(payload.Data) != 1 {
		t.Fatalf("expected 1 visitor in list, got %d", len(payload.Data))
	}
}

func TestVisitorHandlerStatsAndUpdateDuration(t *testing.T) {
	router := setupVisitorTestRouter(t)

	createRequest := httptest.NewRequest(http.MethodPost, "/api/visitors", strings.NewReader(`{"ip":"112.17.45.201","city":"杭州","device":"Windows","browser":"Chrome","os":"Windows","visitTime":"`+time.Now().UTC().Format(time.RFC3339)+`","duration":0}`))
	createRequest.Header.Set("Content-Type", "application/json")
	createRecorder := httptest.NewRecorder()
	router.ServeHTTP(createRecorder, createRequest)

	updateRecorder := httptest.NewRecorder()
	updateRequest := httptest.NewRequest(http.MethodPatch, "/api/visitors/1", strings.NewReader(`{"duration":180}`))
	updateRequest.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(updateRecorder, updateRequest)

	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 from update, got %d", updateRecorder.Code)
	}

	statsRecorder := httptest.NewRecorder()
	statsRequest := httptest.NewRequest(http.MethodGet, "/api/visitors/stats?days=7", nil)
	router.ServeHTTP(statsRecorder, statsRequest)

	if statsRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 from stats, got %d", statsRecorder.Code)
	}

	var statsPayload struct {
		Data struct {
			TotalVisits            int `json:"totalVisits"`
			AverageDurationSeconds int `json:"averageDurationSeconds"`
		} `json:"data"`
	}
	if err := json.Unmarshal(statsRecorder.Body.Bytes(), &statsPayload); err != nil {
		t.Fatalf("expected stats JSON, got error: %v", err)
	}

	if statsPayload.Data.TotalVisits != 1 {
		t.Fatalf("expected total visits 1, got %d", statsPayload.Data.TotalVisits)
	}

	if statsPayload.Data.AverageDurationSeconds != 180 {
		t.Fatalf("expected average duration 180, got %d", statsPayload.Data.AverageDurationSeconds)
	}
}
