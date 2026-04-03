package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"jianli/internal/httpapi"
	"jianli/internal/model"
	"jianli/internal/store"
)

type VisitorHandler struct {
	store *store.VisitorStore
}

type createVisitorRequest struct {
	IP        string    `json:"ip"`
	City      string    `json:"city"`
	Device    string    `json:"device"`
	Browser   string    `json:"browser"`
	OS        string    `json:"os"`
	VisitTime time.Time `json:"visitTime"`
	Duration  int       `json:"duration"`
}

type updateVisitorDurationRequest struct {
	Duration int `json:"duration"`
}

func NewVisitorHandler(store *store.VisitorStore) *VisitorHandler {
	return &VisitorHandler{store: store}
}

func (h *VisitorHandler) Create(c *gin.Context) {
	var request createVisitorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpapi.Error(c, http.StatusBadRequest, 40000, "invalid request body")
		return
	}

	if request.VisitTime.IsZero() {
		request.VisitTime = time.Now().UTC()
	}

	id, err := h.store.RecordVisit(model.VisitorRecord{
		IP:        request.IP,
		City:      request.City,
		Device:    request.Device,
		Browser:   request.Browser,
		OS:        request.OS,
		VisitTime: request.VisitTime,
		Duration:  request.Duration,
	})
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to record visitor")
		return
	}

	httpapi.JSON(c, http.StatusOK, gin.H{"id": id})
}

func (h *VisitorHandler) List(c *gin.Context) {
	days := queryInt(c, "days", 7)
	page := queryInt(c, "page", 1)
	limit := queryInt(c, "limit", 20)

	records, err := h.store.List(days, page, limit)
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to list visitors")
		return
	}

	httpapi.JSON(c, http.StatusOK, records)
}

func (h *VisitorHandler) Stats(c *gin.Context) {
	days := queryInt(c, "days", 7)

	stats, err := h.store.Stats(days)
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to load visitor stats")
		return
	}

	httpapi.JSON(c, http.StatusOK, stats)
}

func (h *VisitorHandler) UpdateDuration(c *gin.Context) {
	visitorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.Error(c, http.StatusBadRequest, 40000, "invalid visitor id")
		return
	}

	var request updateVisitorDurationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpapi.Error(c, http.StatusBadRequest, 40000, "invalid request body")
		return
	}

	if err := h.store.UpdateDuration(visitorID, request.Duration); err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to update visitor duration")
		return
	}

	httpapi.JSON(c, http.StatusOK, gin.H{"updated": true})
}

func queryInt(c *gin.Context, key string, fallback int) int {
	value := c.Query(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
