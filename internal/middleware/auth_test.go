package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddlewareAllowsMatchingKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(Auth("resume-key"))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("X-Auth-Key", "resume-key")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", recorder.Code)
	}
}

func TestAuthMiddlewareRejectsInvalidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(Auth("resume-key"))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("X-Auth-Key", "wrong-key")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}
