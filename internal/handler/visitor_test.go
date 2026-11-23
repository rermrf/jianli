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
	"jianli/internal/model"
	"jianli/internal/store"
)

type visitorLocationResolver struct {
	calls []string
}

func (r *visitorLocationResolver) Resolve(ip string) (country, region, city, isp string, ok bool) {
	r.calls = append(r.calls, ip)
	if ip == "112.17.45.201" || ip == "112.17.45.202" {
		return "中国", "浙江", "杭州", "移动", true
	}
	if ip == "183.62.84.19" {
		return "中国", "广东", "深圳", "电信", true
	}
	return "", "", "", "", false
}

func TestVisitorHandlerListRepairsAndPersistsUnknownLocations(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	visitorStore := store.NewVisitorStore(db)
	unknown := model.VisitorRecord{IP: "112.17.45.201", City: "未知", VisitTime: time.Now(), Duration: 120, PDFExported: true}
	id, err := visitorStore.RecordVisit(unknown)
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range []model.VisitorRecord{
		{IP: "183.62.84.19", Country: "中国", Region: "广东", City: "深圳", VisitTime: time.Now()},
		{IP: "112.17.45.202", Country: "中国", Region: "0", City: "未知", VisitTime: time.Now()},
		{IP: "invalid-ip", City: "未知", VisitTime: time.Now()},
	} {
		if _, err := visitorStore.RecordVisit(record); err != nil {
			t.Fatal(err)
		}
	}
	resolver := &visitorLocationResolver{}
	h := NewVisitorHandler(visitorStore, resolver)
	router := gin.New()
	router.GET("/api/visitors", h.List)
	for i := 0; i < 2; i++ {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/visitors?days=0", nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("list status=%d body=%s", recorder.Code, recorder.Body.String())
		}
		var payload struct {
			Data []model.VisitorRecord `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatal(err)
		}
		if len(payload.Data) != 4 {
			t.Fatalf("expected all 4 existing rows, got %d", len(payload.Data))
		}
		for _, record := range payload.Data {
			if record.ID == id && (record.Country != "中国" || record.Region != "浙江" || record.City != "杭州" || record.ISP != "移动") {
				t.Fatalf("list did not return repaired geography: %+v", record)
			}
			if record.IP == "112.17.45.202" && (record.Country != "中国" || record.Region != "浙江" || record.City != "杭州" || record.ISP != "移动") {
				t.Fatalf("country-only geography was not completed: %+v", record)
			}
			if record.IP == "183.62.84.19" && (record.Country != "中国" || record.Region != "广东" || record.City != "深圳" || record.ISP != "电信") {
				t.Fatalf("known location must be preserved while filling ISP: %+v", record)
			}
		}
	}
	var saved model.VisitorRecord
	if err := db.First(&saved, id).Error; err != nil {
		t.Fatal(err)
	}
	if saved.City != "杭州" || saved.Country != "中国" || saved.Region != "浙江" || saved.ISP != "移动" || saved.Duration != 120 || !saved.PDFExported || saved.IP != unknown.IP {
		t.Fatalf("unexpected persisted row: %+v", saved)
	}
	resolvedCalls := 0
	for _, ip := range resolver.calls {
		if ip == unknown.IP {
			resolvedCalls++
		}
	}
	if resolvedCalls != 1 {
		t.Fatalf("repaired rows should not be resolved again, got %d calls", resolvedCalls)
	}
}

func TestVisitorHandlerRecordsDetailedDeviceMetadata(t *testing.T) {
	for _, tc := range []struct {
		name, ua, device, os, osVersion, browser, browserVersion string
	}{
		{"Windows Edge", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/146.0.0.0 Safari/537.36 Edg/146.0.3856.84", "Desktop", "Windows", "10/11", "Edge", "146.0.3856.84"},
		{"Android Edge", "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 Chrome/146.0.0.0 Mobile Safari/537.36 EdgA/146.0.0.0", "Mobile", "Android", "14", "Edge", "146.0.0.0"},
		{"Android tablet", "Mozilla/5.0 (Linux; Android 13; SM-X700) AppleWebKit/537.36 Chrome/146.0.0.0 Safari/537.36", "Tablet", "Android", "13", "Chrome", "146.0.0.0"},
		{"iOS Chrome", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_5 like Mac OS X) AppleWebKit/605.1.15 CriOS/146.0.7680.90 Mobile/15E148 Safari/604.1", "Mobile", "iOS", "18.5", "Chrome", "146.0.7680.90"},
		{"iOS Firefox", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_5 like Mac OS X) AppleWebKit/605.1.15 FxiOS/146.0 Mobile/15E148 Safari/605.1.15", "Mobile", "iOS", "18.5", "Firefox", "146.0"},
		{"iPad Safari", "Mozilla/5.0 (iPad; CPU OS 18_5 like Mac OS X) AppleWebKit/605.1.15 Version/18.5 Mobile/15E148 Safari/604.1", "Tablet", "iOS", "18.5", "Safari", "18.5"},
		{"iPad desktop mode", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15) AppleWebKit/605.1.15 Version/18.5 Mobile/15E148 Safari/604.1", "Tablet", "iOS", "", "Safari", "18.5"},
		{"macOS Safari", "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/605.1.15 Version/17.5 Safari/605.1.15", "Desktop", "macOS", "14.5", "Safari", "17.5"},
		{"Linux Firefox", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:146.0) Gecko/20100101 Firefox/146.0", "Desktop", "Linux", "", "Firefox", "146.0"},
		{"bot", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", "Bot", "Unknown", "", "Unknown", ""},
		{"missing user agent", "", "Unknown", "Unknown", "", "Unknown", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			router := setupVisitorTestRouter(t, geoip.NewFakeResolver("中国", "浙江省", "杭州市", "移动", true))
			request := httptest.NewRequest(http.MethodPost, "/api/visitors", strings.NewReader(`{"duration":12,"ip":"39.182.***.***","os":"fake"}`))
			request.RemoteAddr = "39.182.123.234:54321"
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("User-Agent", tc.ua)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusOK {
				t.Fatalf("create status=%d body=%s", recorder.Code, recorder.Body.String())
			}
			list := httptest.NewRecorder()
			router.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/api/visitors", nil))
			var payload struct {
				Data []model.VisitorRecord `json:"data"`
			}
			if err := json.Unmarshal(list.Body.Bytes(), &payload); err != nil {
				t.Fatal(err)
			}
			if len(payload.Data) != 1 {
				t.Fatalf("unexpected visitor count: %d", len(payload.Data))
			}
			got := payload.Data[0]
			if got.Device != tc.device || got.OS != tc.os || got.OSVersion != tc.osVersion || got.Browser != tc.browser || got.BrowserVersion != tc.browserVersion || got.UserAgent != tc.ua {
				t.Fatalf("incorrect recorded metadata: %+v", got)
			}
			if got.IP != "39.182.123.234" || got.Country != "中国" || got.Region != "浙江省" || got.City != "杭州市" || got.ISP != "移动" {
				t.Fatalf("expected full IP and location: %+v", got)
			}
		})
	}
}

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
		Code int              `json:"code"`
		Data []map[string]any `json:"data"`
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

func TestVisitorHandlerTrendRejectsNegativeDays(t *testing.T) {
	router := setupVisitorTestRouter(t, nil)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/visitors/trend?days=-1", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for days=-1, got %d", recorder.Code)
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
		Data []model.VisitorTrendPoint `json:"data"`
	}
	if err := json.Unmarshal(trendRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("trend JSON: %v", err)
	}
	if len(payload.Data) != 7 {
		t.Fatalf("expected one bucket per day, got %d", len(payload.Data))
	}
	today := payload.Data[6]
	if today.Count != 1 || today.UniqueVisitors != 1 || len(today.Date) != 10 {
		t.Fatalf("expected today's visitor and canonical date, got %+v", today)
	}
}

func TestVisitorHandlerTrendSupportsAllHistory(t *testing.T) {
	router := setupVisitorTestRouter(t, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/visitors/trend?days=0", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("all-history status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var payload struct {
		Data []model.VisitorTrendPoint `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Data == nil || len(payload.Data) != 0 {
		t.Fatalf("expected an empty JSON array, got %#v", payload.Data)
	}
}
