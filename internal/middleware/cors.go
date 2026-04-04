package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Access-Control-Allow-Headers", "Content-Type, X-Auth-Key")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
