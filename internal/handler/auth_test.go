package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestVerifyHandlerReturnsValidTrueForMatchingKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/api/auth/verify", VerifyAuth("resume-key"))

	request := httptest.NewRequest(http.MethodPost, "/api/auth/verify", strings.NewReader(`{"key":"resume-key"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var payload struct {
		Code int `json:"code"`
		Data struct {
			Valid bool `json:"valid"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}

	if !payload.Data.Valid {
		t.Fatal("expected valid=true")
	}
}

func TestVerifyHandlerRejectsInvalidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/api/auth/verify", VerifyAuth("resume-key"))

	request := httptest.NewRequest(http.MethodPost, "/api/auth/verify", strings.NewReader(`{"key":"wrong-key"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}
