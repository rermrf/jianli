package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"jianli/internal/middleware"
)

func buildMultipartRequest(t *testing.T, fieldName, fileName string, content []byte) (*bytes.Buffer, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("CreateFormFile() returned error: %v", err)
	}

	if _, err := part.Write(content); err != nil {
		t.Fatalf("part.Write() returned error: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() returned error: %v", err)
	}

	return body, writer.FormDataContentType()
}

func setupUploadTestRouter(t *testing.T, uploadDir string) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	protected := router.Group("/api")
	protected.Use(middleware.Auth("resume-key"))
	protected.POST("/upload/avatar", NewUploadHandler(uploadDir).UploadAvatar)

	return router
}

func TestUploadAvatarRequiresAuth(t *testing.T) {
	router := setupUploadTestRouter(t, filepath.Join(t.TempDir(), "avatars"))

	body, contentType := buildMultipartRequest(t, "file", "avatar.png", []byte{0x89, 0x50, 0x4E, 0x47})
	request := httptest.NewRequest(http.MethodPost, "/api/upload/avatar", body)
	request.Header.Set("Content-Type", contentType)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestUploadAvatarRejectsNonImageFiles(t *testing.T) {
	router := setupUploadTestRouter(t, filepath.Join(t.TempDir(), "avatars"))

	body, contentType := buildMultipartRequest(t, "file", "avatar.txt", []byte("plain-text"))
	request := httptest.NewRequest(http.MethodPost, "/api/upload/avatar", body)
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("X-Auth-Key", "resume-key")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestUploadAvatarStoresImageAndReturnsURL(t *testing.T) {
	uploadDir := filepath.Join(t.TempDir(), "avatars")
	router := setupUploadTestRouter(t, uploadDir)

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

	body, contentType := buildMultipartRequest(t, "file", "avatar.png", pngBytes)
	request := httptest.NewRequest(http.MethodPost, "/api/upload/avatar", body)
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("X-Auth-Key", "resume-key")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var payload struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}

	if !strings.HasPrefix(payload.Data.URL, "/uploads/avatars/") {
		t.Fatalf("expected upload url prefix, got %q", payload.Data.URL)
	}

	savedPath := filepath.Join(uploadDir, filepath.Base(payload.Data.URL))
	if _, err := os.Stat(savedPath); err != nil {
		t.Fatalf("expected uploaded avatar to be stored at %q: %v", savedPath, err)
	}
}
