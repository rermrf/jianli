package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"jianli/internal/httpapi"
	"jianli/internal/model"
	"jianli/internal/store"
)

type ResumeDraftHandler struct {
	store *store.ResumeDraftStore
}

type createResumeDraftRequest struct {
	Name string          `json:"name"`
	Note string          `json:"note"`
	Data json.RawMessage `json:"data"`
}

type resumeDraftSummaryResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Note      string `json:"note"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type resumeDraftDetailResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Note      string `json:"note"`
	Data      any    `json:"data"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func NewResumeDraftHandler(store *store.ResumeDraftStore) *ResumeDraftHandler {
	return &ResumeDraftHandler{store: store}
}

func (h *ResumeDraftHandler) Create(c *gin.Context) {
	var request createResumeDraftRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpapi.Error(c, http.StatusBadRequest, 40000, "invalid request body")
		return
	}
	if strings.TrimSpace(request.Name) == "" || len(request.Data) == 0 {
		httpapi.Error(c, http.StatusBadRequest, 40000, "invalid request body")
		return
	}

	id, err := h.store.Create(strings.TrimSpace(request.Name), request.Note, request.Data)
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to save resume draft")
		return
	}

	record, err := h.store.Get(id)
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to load resume draft")
		return
	}

	response, err := toResumeDraftDetailResponse(record)
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to decode resume draft")
		return
	}

	httpapi.JSON(c, http.StatusOK, response)
}

func (h *ResumeDraftHandler) List(c *gin.Context) {
	records, err := h.store.List()
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to list resume drafts")
		return
	}

	response := make([]resumeDraftSummaryResponse, 0, len(records))
	for _, record := range records {
		response = append(response, toResumeDraftSummaryResponse(record))
	}

	httpapi.JSON(c, http.StatusOK, response)
}

func (h *ResumeDraftHandler) Get(c *gin.Context) {
	id, ok := parseResumeDraftID(c)
	if !ok {
		return
	}

	record, err := h.store.Get(id)
	if err != nil {
		writeResumeDraftStoreError(c, err, "failed to load resume draft")
		return
	}

	response, err := toResumeDraftDetailResponse(record)
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to decode resume draft")
		return
	}

	httpapi.JSON(c, http.StatusOK, response)
}

func (h *ResumeDraftHandler) Publish(c *gin.Context) {
	id, ok := parseResumeDraftID(c)
	if !ok {
		return
	}

	if err := h.store.Publish(id); err != nil {
		writeResumeDraftStoreError(c, err, "failed to publish resume draft")
		return
	}

	httpapi.JSON(c, http.StatusOK, gin.H{"published": true})
}

func (h *ResumeDraftHandler) Delete(c *gin.Context) {
	id, ok := parseResumeDraftID(c)
	if !ok {
		return
	}

	if err := h.store.Delete(id); err != nil {
		writeResumeDraftStoreError(c, err, "failed to delete resume draft")
		return
	}

	httpapi.JSON(c, http.StatusOK, gin.H{"deleted": true})
}

func parseResumeDraftID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.Error(c, http.StatusBadRequest, 40000, "invalid draft id")
		return 0, false
	}

	return id, true
}

func writeResumeDraftStoreError(c *gin.Context, err error, fallbackMessage string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpapi.Error(c, http.StatusNotFound, 40404, "resume draft not found")
		return
	}

	httpapi.Error(c, http.StatusInternalServerError, 50000, fallbackMessage)
}

func toResumeDraftSummaryResponse(record model.ResumeDraftRecord) resumeDraftSummaryResponse {
	return resumeDraftSummaryResponse{
		ID:        record.ID,
		Name:      record.Name,
		Note:      record.Note,
		CreatedAt: record.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: record.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toResumeDraftDetailResponse(record model.ResumeDraftRecord) (resumeDraftDetailResponse, error) {
	var payload any
	if err := json.Unmarshal(record.Data, &payload); err != nil {
		return resumeDraftDetailResponse{}, err
	}

	return resumeDraftDetailResponse{
		ID:        record.ID,
		Name:      record.Name,
		Note:      record.Note,
		Data:      payload,
		CreatedAt: record.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: record.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}
