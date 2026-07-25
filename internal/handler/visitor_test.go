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

	"jianli/internal/geoip"
	"jianli/internal/store"
)

func setupVisitorTestRouter(t *testing.T, resolver geoip.Resolver) *gin.Engine {
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

	if resolver == nil {
		resolver = geoip.NewLoopbackResolver()
	}
	visitorHandler := NewVisitorHandler(store.NewVisitorStore(db), resolver)
	router := gin.New()
	router.POST("/api/visitors", visitorHandler.Create)
	router.POST("/api/visitors/:id", visitorHandler.UpdateDuration)
	router.PATCH("/api/visitors/:id", visitorHandler.UpdateDuration)
	router.POST("/api/visitors/:id/pdf-export", visitorHandler.MarkPDFExported)
	router.GET("/api/visitors", visitorHandler.List)
	router.GET("/api/visitors/stats", visitorHandler.Stats)
	router.GET("/api/visitors/trend", visitorHandler.Trend)

	return router
}

// nowRFC3339 returns the current UTC time formatted as RFC3339. Tests use
// this instead of hardcoded dates so the days-window filter on the List
// endpoint always includes the freshly-recorded visit (a hardcoded date
// older than 7 days would silently disappear from the list).
func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func TestVisitorHandlerCreateAndList(t *testing.T) {
	router := setupVisitorTestRouter(t, nil)

	body := `{"visitTime":"` + nowRFC3339() + `","duration":0}`
	request := httptest.NewRequest(http.MethodPost, "/api/visitors", strings.NewReader(body))
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
		Code int                `json:"code"`
		Data []map[string]any   `json:"data"`
	}
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected list JSON, got error: %v", err)
	}

	if len(payload.Data) != 1 {
		t.Fatalf("expected 1 visitor in list, got %d", len(payload.Data))
	}
}

func TestVisitorHandlerCreateDerivesMetadataFromRequest(t *testing.T) {
	router := setupVisitorTestRouter(t, nil)

	body := `{"visitTime":"` + nowRFC3339() + `","duration":0}`
	request := httptest.NewRequest(http.MethodPost, "/api/visitors", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36 Edg/146.0.0.0")
	request.RemoteAddr = "127.0.0.1:54321"
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 from create, got %d", recorder.Code)
	}

	listRecorder := httptest.NewRecorder()
	listRequest := httptest.NewRequest(http.MethodGet, "/api/visitors?days=7&page=1&limit=20", nil)
	router.ServeHTTP(listRecorder, listRequest)

	var payload struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected list JSON, got error: %v", err)
	}

	if len(payload.Data) == 0 {
		t.Fatalf("expected at least one record, got 0")
	}

	record := payload.Data[0]
	if record["city"] != "本地网络" {
		t.Fatalf("expected derived city 本地网络, got %#v", record["city"])
	}
	if record["country"] != "本地" {
		t.Fatalf("expected derived country 本地, got %#v", record["country"])
	}
	if record["browser"] != "Edge" {
		t.Fatalf("expected derived browser Edge, got %#v", record["browser"])
	}
	if record["device"] != "Desktop" {
		t.Fatalf("expected derived device Desktop, got %#v", record["device"])
	}
}

func TestVisitorHandlerCreateUsesGeoLookup(t *testing.T) {
	resolver := geoip.NewFakeResolver("中国", "浙江", "杭州", "电信", true)
	router := setupVisitorTestRouter(t, resolver)

	// Use a public-looking IP that is neither loopback nor private; the
	// fake resolver returns the canned values regardless.
	body := `{"visitTime":"` + nowRFC3339() + `","duration":0}`
	request := httptest.NewRequest(http.MethodPost, "/api/visitors", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.RemoteAddr = "112.17.45.201:54321"
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	listRecorder := httptest.NewRecorder()
	listRequest := httptest.NewRequest(http.MethodGet, "/api/visitors?days=7", nil)
	router.ServeHTTP(listRecorder, listRequest)

	var payload struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("list JSON: %v", err)
	}
	if len(payload.Data) == 0 {
		t.Fatal("expected one record")
	}
	record := payload.Data[0]
	if record["country"] != "中国" || record["region"] != "浙江" || record["city"] != "杭州" || record["isp"] != "电信" {
		t.Fatalf("expected resolver values, got %#v", record)
	}
}

func TestVisitorHandlerStatsAndUpdateDuration(t *testing.T) {
	router := setupVisitorTestRouter(t, nil)

	createRequest := httptest.NewRequest(http.MethodPost, "/api/visitors", strings.NewReader(`{"visitTime":"`+nowRFC3339()+`","duration":0}`))
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

func TestVisitorHandlerMarkPDFExportedFlipsFlag(t *testing.T) {
	router := setupVisitorTestRouter(t, nil)

	createRecorder := httptest.NewRecorder()
	createRequest := httptest.NewRequest(http.MethodPost, "/api/visitors", strings.NewReader(`{"visitTime":"`+nowRFC3339()+`","duration":0}`))
	createRequest.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(createRecorder, createRequest)

	markRecorder := httptest.NewRecorder()
	markRequest := httptest.NewRequest(http.MethodPost, "/api/visitors/1/pdf-export", nil)
	router.ServeHTTP(markRecorder, markRequest)

	if markRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 from pdf-export, got %d", markRecorder.Code)
	}

	listRecorder := httptest.NewRecorder()
	listRequest := httptest.NewRequest(http.MethodGet, "/api/visitors?days=7", nil)
	router.ServeHTTP(listRecorder, listRequest)

	var payload struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("list JSON: %v", err)
	}
	if got := payload.Data[0]["pdfExported"]; got != true {
		t.Fatalf("expected pdfExported=true, got %#v", got)
	}
}

func TestVisitorHandlerTrendRejectsZeroDays(t *testing.T) {
	router := setupVisitorTestRouter(t, nil)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/visitors/trend?days=0", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for days=0, got %d", recorder.Code)
	}
}

func TestVisitorHandlerTrendReturnsBucketedSeries(t *testing.T) {
	router := setupVisitorTestRouter(t, nil)

	createRequest := httptest.NewRequest(http.MethodPost, "/api/visitors", strings.NewReader(`{"visitTime":"`+nowRFC3339()+`","duration":0}`))
	createRequest.Header.Set("Content-Type", "application/json")
	createRecorder := httptest.NewRecorder()
	router.ServeHTTP(createRecorder, createRequest)

	trendRecorder := httptest.NewRecorder()
	trendRequest := httptest.NewRequest(http.MethodGet, "/api/visitors/trend?days=7", nil)
	router.ServeHTTP(trendRecorder, trendRequest)

	if trendRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", trendRecorder.Code)
	}

	var payload struct {
		Data []struct {
			Date  string `json:"date"`
			Count int    `json:"count"`
		} `json:"data"`
	}
	if err := json.Unmarshal(trendRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("trend JSON: %v", err)
	}
	if len(payload.Data) == 0 {
		t.Fatal("expected at least one bucket for today's visit")
	}
	if payload.Data[0].Count < 1 {
		t.Fatalf("expected count >= 1, got %d", payload.Data[0].Count)
	}
}
