package handler

import (
	"net"
	"net/http"
	"strconv"
	"strings"
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

	visitTime := request.VisitTime
	if visitTime.IsZero() {
		visitTime = time.Now().UTC()
	}

	userAgent := c.GetHeader("User-Agent")
	browser, osName, device := deriveUserAgentMetadata(userAgent)
	ip := c.ClientIP()
	city := deriveCity(ip)

	id, err := h.store.RecordVisit(model.VisitorRecord{
		IP:        ip,
		City:      city,
		Device:    device,
		Browser:   browser,
		OS:        osName,
		VisitTime: visitTime,
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

func deriveCity(ip string) string {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "未知"
	}
	if parsedIP.IsLoopback() || parsedIP.IsPrivate() {
		return "本地网络"
	}
	return "未知"
}

func deriveUserAgentMetadata(userAgent string) (browser string, osName string, device string) {
	lowerUA := strings.ToLower(userAgent)

	switch {
	case strings.Contains(lowerUA, "edg/"):
		browser = "Edge"
	case strings.Contains(lowerUA, "chrome/"):
		browser = "Chrome"
	case strings.Contains(lowerUA, "firefox/"):
		browser = "Firefox"
	case strings.Contains(lowerUA, "safari/") && !strings.Contains(lowerUA, "chrome/"):
		browser = "Safari"
	default:
		browser = "Unknown"
	}

	switch {
	case strings.Contains(lowerUA, "windows"):
		osName = "Windows"
	case strings.Contains(lowerUA, "android"):
		osName = "Android"
	case strings.Contains(lowerUA, "iphone") || strings.Contains(lowerUA, "ipad") || strings.Contains(lowerUA, "ios"):
		osName = "iOS"
	case strings.Contains(lowerUA, "mac os x") || strings.Contains(lowerUA, "macintosh"):
		osName = "macOS"
	case strings.Contains(lowerUA, "linux"):
		osName = "Linux"
	default:
		osName = "Unknown"
	}

	switch {
	case strings.Contains(lowerUA, "iphone"), strings.Contains(lowerUA, "android") && strings.Contains(lowerUA, "mobile"):
		device = "Mobile"
	case strings.Contains(lowerUA, "ipad"), strings.Contains(lowerUA, "tablet"):
		device = "Tablet"
	default:
		device = "Desktop"
	}

	return browser, osName, device
}
