package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jianli/internal/httpapi"
)

type verifyAuthRequest struct {
	Key string `json:"key"`
}

func VerifyAuth(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request verifyAuthRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			httpapi.Error(c, http.StatusBadRequest, 40000, "invalid request body")
			return
		}

		if request.Key != expectedKey {
			httpapi.Error(c, http.StatusUnauthorized, 40001, "invalid auth key")
			return
		}

		httpapi.JSON(c, http.StatusOK, gin.H{"valid": true})
	}
}
