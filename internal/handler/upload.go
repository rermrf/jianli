package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"jianli/internal/httpapi"
)

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	return &UploadHandler{uploadDir: uploadDir}
}

func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		httpapi.Error(c, http.StatusBadRequest, 40000, "file is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		httpapi.Error(c, http.StatusBadRequest, 40000, "failed to open uploaded file")
		return
	}
	defer file.Close()

	content := make([]byte, fileHeader.Size)
	if _, err := file.Read(content); err != nil {
		httpapi.Error(c, http.StatusBadRequest, 40000, "failed to read uploaded file")
		return
	}

	contentType := http.DetectContentType(content)
	if !strings.HasPrefix(contentType, "image/") {
		httpapi.Error(c, http.StatusBadRequest, 40000, "file must be an image")
		return
	}

	if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to prepare upload directory")
		return
	}

	ext := extensionFromContentType(contentType)
	fileName := fmt.Sprintf("avatar-%d%s", time.Now().UnixNano(), ext)
	savedPath := filepath.Join(h.uploadDir, fileName)
	if err := os.WriteFile(savedPath, content, 0o644); err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to store avatar")
		return
	}

	httpapi.JSON(c, http.StatusOK, gin.H{
		"url": "/uploads/avatars/" + fileName,
	})
}

func extensionFromContentType(contentType string) string {
	switch contentType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	default:
		return ".bin"
	}
}
