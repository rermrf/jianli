package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"jianli/internal/httpapi"
	"jianli/internal/store"
)

type PDFExporter interface {
	ExportResume(ctx context.Context, resume json.RawMessage) ([]byte, error)
}

type ResumeHandler struct {
	exporter PDFExporter
	store    *store.ResumeStore
}

func NewResumeHandler(store *store.ResumeStore, exporter PDFExporter) *ResumeHandler {
	return &ResumeHandler{exporter: exporter, store: store}
}

func (h *ResumeHandler) Get(c *gin.Context) {
	resume, err := h.store.Get()
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to load resume")
		return
	}

	var payload any
	if err := json.Unmarshal(resume, &payload); err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to decode resume")
		return
	}

	httpapi.JSON(c, http.StatusOK, payload)
}

func (h *ResumeHandler) Update(c *gin.Context) {
	var payload map[string]any
	if err := c.ShouldBindJSON(&payload); err != nil {
		httpapi.Error(c, http.StatusBadRequest, 40000, "invalid request body")
		return
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to encode resume")
		return
	}

	if err := h.store.Save(encoded); err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to save resume")
		return
	}

	httpapi.JSON(c, http.StatusOK, payload)
}

func (h *ResumeHandler) ExportPDF(c *gin.Context) {
	resume, err := h.store.Get()
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to load resume")
		return
	}

	pdfBytes, err := h.exporter.ExportResume(c.Request.Context(), resume)
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to export pdf")
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=resume.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
