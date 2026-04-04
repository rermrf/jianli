package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jianli/internal/httpapi"
)

func Auth(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Auth-Key") != expectedKey {
			httpapi.Error(c, http.StatusUnauthorized, 40001, "invalid auth key")
			c.Abort()
			return
		}

		c.Next()
	}
}
