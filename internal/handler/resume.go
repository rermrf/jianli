package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"jianli/internal/httpapi"
	"jianli/internal/store"
)

type ResumeHandler struct {
	store *store.ResumeStore
}

func NewResumeHandler(store *store.ResumeStore) *ResumeHandler {
	return &ResumeHandler{store: store}
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
