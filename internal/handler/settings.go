package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jianli/internal/httpapi"
	"jianli/internal/store"
)

type SettingsHandler struct {
	store *store.SiteSettingsStore
}

type settingsPayload struct {
	AllowPDFExport bool `json:"allowPdfExport"`
}

func NewSettingsHandler(store *store.SiteSettingsStore) *SettingsHandler {
	return &SettingsHandler{store: store}
}

func (h *SettingsHandler) Get(c *gin.Context) {
	settings, err := h.store.Get()
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to load settings")
		return
	}

	httpapi.JSON(c, http.StatusOK, settingsPayload{AllowPDFExport: settings.AllowPDFExport})
}

func (h *SettingsHandler) Update(c *gin.Context) {
	var payload settingsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		httpapi.Error(c, http.StatusBadRequest, 40000, "invalid request body")
		return
	}

	if err := h.store.Save(payload.AllowPDFExport); err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to save settings")
		return
	}

	httpapi.JSON(c, http.StatusOK, payload)
}
