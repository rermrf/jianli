package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"jianli/internal/middleware"
	"jianli/internal/store"
)

type fakeDraftPDFExporter struct {
	content []byte
}

func (f fakeDraftPDFExporter) ExportResume(_ context.Context, _ json.RawMessage) ([]byte, error) {
	return f.content, nil
}

func setupResumeDraftTestRouter(t *testing.T) *gin.Engine {
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

	resumeStore := store.NewResumeStore(db)
	resumeHandler := NewResumeHandler(resumeStore, fakeDraftPDFExporter{content: []byte("%PDF-test")})
	draftHandler := NewResumeDraftHandler(store.NewResumeDraftStore(db, resumeStore))

	router := gin.New()
	router.GET("/api/resume", resumeHandler.Get)
	protected := router.Group("/api")
	protected.Use(middleware.Auth("resume-key"))
	protected.PUT("/resume", resumeHandler.Update)
	protected.POST("/resume/drafts", draftHandler.Create)
	protected.GET("/resume/drafts", draftHandler.List)
	protected.GET("/resume/drafts/:id", draftHandler.Get)
	protected.PUT("/resume/drafts/:id/publish", draftHandler.Publish)
	protected.DELETE("/resume/drafts/:id", draftHandler.Delete)

	return router
}

func createResumeDraftThroughAPI(t *testing.T, router *gin.Engine, body string) int64 {
	t.Helper()

	request := httptest.NewRequest(http.MethodPost, "/api/resume/drafts", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Key", "resume-key")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected create draft status 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected create draft response JSON, got error: %v", err)
	}

	return response.Data.ID
}

func TestResumeDraftHandlerCreateRequiresAuth(t *testing.T) {
	router := setupResumeDraftTestRouter(t)
	request := httptest.NewRequest(http.MethodPost, "/api/resume/drafts", strings.NewReader(`{"name":"未授权草稿","note":"不应创建","data":{"profile":{"name":"未授权"}}}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestResumeDraftHandlerCreateRejectsEmptyName(t *testing.T) {
	router := setupResumeDraftTestRouter(t)
	request := httptest.NewRequest(http.MethodPost, "/api/resume/drafts", strings.NewReader(`{"name":"","note":"缺少名称","data":{"profile":{"name":"无效草稿"}}}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Key", "resume-key")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d with body %s", recorder.Code, recorder.Body.String())
	}
}

func TestResumeDraftHandlerListAndGetReturnDraftData(t *testing.T) {
	router := setupResumeDraftTestRouter(t)
	firstID := createResumeDraftThroughAPI(t, router, `{"name":"初版草稿","note":"最早的草稿","data":{"profile":{"name":"初版"}}}`)
	secondID := createResumeDraftThroughAPI(t, router, `{"name":"终版草稿","note":"最近一次修改","data":{"profile":{"name":"终版"}}}`)

	listRequest := httptest.NewRequest(http.MethodGet, "/api/resume/drafts", nil)
	listRequest.Header.Set("X-Auth-Key", "resume-key")
	listRecorder := httptest.NewRecorder()
	router.ServeHTTP(listRecorder, listRequest)

	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", listRecorder.Code, listRecorder.Body.String())
	}

	var listResponse struct {
		Data []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Note string `json:"note"`
		} `json:"data"`
	}
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &listResponse); err != nil {
		t.Fatalf("expected list response JSON, got error: %v", err)
	}

	if len(listResponse.Data) != 2 {
		t.Fatalf("expected 2 drafts, got %d", len(listResponse.Data))
	}
	if listResponse.Data[0].ID != secondID || listResponse.Data[1].ID != firstID {
		t.Fatalf("expected drafts ordered newest first, got %#v", listResponse.Data)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/api/resume/drafts/"+strconv.FormatInt(secondID, 10), nil)
	getRequest.Header.Set("X-Auth-Key", "resume-key")
	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", getRecorder.Code, getRecorder.Body.String())
	}

	var getResponse struct {
		Data struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Note string `json:"note"`
			Data struct {
				Profile struct {
					Name string `json:"name"`
				} `json:"profile"`
			} `json:"data"`
		} `json:"data"`
	}
	if err := json.Unmarshal(getRecorder.Body.Bytes(), &getResponse); err != nil {
		t.Fatalf("expected get response JSON, got error: %v", err)
	}

	if getResponse.Data.ID != secondID || getResponse.Data.Name != "终版草稿" || getResponse.Data.Note != "最近一次修改" {
		t.Fatalf("expected latest draft detail, got %#v", getResponse.Data)
	}
	if getResponse.Data.Data.Profile.Name != "终版" {
		t.Fatalf("expected detail payload to include draft resume data, got %#v", getResponse.Data.Data.Profile.Name)
	}
}

func TestResumeDraftHandlerPublishUpdatesMainResume(t *testing.T) {
	router := setupResumeDraftTestRouter(t)
	draftID := createResumeDraftThroughAPI(t, router, `{"name":"发布版","note":"准备用于线上","data":{"profile":{"name":"已发布简历"}}}`)

	publishRequest := httptest.NewRequest(http.MethodPut, "/api/resume/drafts/"+strconv.FormatInt(draftID, 10)+"/publish", nil)
	publishRequest.Header.Set("X-Auth-Key", "resume-key")
	publishRecorder := httptest.NewRecorder()
	router.ServeHTTP(publishRecorder, publishRequest)

	if publishRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", publishRecorder.Code, publishRecorder.Body.String())
	}

	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, httptest.NewRequest(http.MethodGet, "/api/resume", nil))

	var response struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(getRecorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected resume response JSON, got error: %v", err)
	}

	profile := response.Data["profile"].(map[string]any)
	if profile["name"] != "已发布简历" {
		t.Fatalf("expected published profile name, got %#v for draft %d", profile["name"], draftID)
	}
}

func TestResumeDraftHandlerDeleteRemovesDraft(t *testing.T) {
	router := setupResumeDraftTestRouter(t)
	draftID := createResumeDraftThroughAPI(t, router, `{"name":"待删除草稿","note":"删除后不应存在","data":{"profile":{"name":"待删除"}}}`)

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/resume/drafts/"+strconv.FormatInt(draftID, 10), nil)
	deleteRequest.Header.Set("X-Auth-Key", "resume-key")
	deleteRecorder := httptest.NewRecorder()
	router.ServeHTTP(deleteRecorder, deleteRequest)

	if deleteRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", deleteRecorder.Code, deleteRecorder.Body.String())
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/api/resume/drafts/"+strconv.FormatInt(draftID, 10), nil)
	getRequest.Header.Set("X-Auth-Key", "resume-key")
	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete for draft %d, got %d with body %s", draftID, getRecorder.Code, getRecorder.Body.String())
	}
}

func TestResumeDraftHandlerMissingDraftReturnsNotFound(t *testing.T) {
	router := setupResumeDraftTestRouter(t)

	for _, request := range []*http.Request{
		httptest.NewRequest(http.MethodGet, "/api/resume/drafts/999", nil),
		httptest.NewRequest(http.MethodPut, "/api/resume/drafts/999/publish", nil),
		httptest.NewRequest(http.MethodDelete, "/api/resume/drafts/999", nil),
	} {
		request.Header.Set("X-Auth-Key", "resume-key")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusNotFound {
			t.Fatalf("expected 404 for %s %s, got %d with body %s", request.Method, request.URL.Path, recorder.Code, recorder.Body.String())
		}
	}
}
