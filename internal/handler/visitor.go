package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"jianli/internal/geoip"
	"jianli/internal/httpapi"
	"jianli/internal/model"
	"jianli/internal/store"
)

type VisitorHandler struct {
	store *store.VisitorStore
	geoip geoip.Resolver
}

type createVisitorRequest struct {
	VisitTime time.Time `json:"visitTime"`
	Duration  int       `json:"duration"`
}

type updateVisitorDurationRequest struct {
	Duration int `json:"duration"`
}

// NewVisitorHandler requires a non-nil Resolver. Callers that fail to load
// the ip2region xdb should pass geoip.NewLoopbackResolver() so loopback /
// private IPs still get "本地网络" while public IPs fall through to "未知".
func NewVisitorHandler(s *store.VisitorStore, resolver geoip.Resolver) *VisitorHandler {
	return &VisitorHandler{store: s, geoip: resolver}
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

	country, region, city, isp, _ := h.geoip.Resolve(ip)
	if city == "" {
		city = "未知"
	}

	id, err := h.store.RecordVisit(model.VisitorRecord{
		IP:        ip,
		Country:   country,
		Region:    region,
		City:      city,
		ISP:       isp,
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

// Trend returns daily-bucketed visit counts. days must be a positive int;
// frontend translates the 'all' range to days=30 to keep buckets reasonable.
func (h *VisitorHandler) Trend(c *gin.Context) {
	days := queryInt(c, "days", 7)
	if days <= 0 {
		httpapi.Error(c, http.StatusBadRequest, 40000, "days must be positive")
		return
	}

	points, err := h.store.Trend(days)
	if err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to load trend")
		return
	}

	httpapi.JSON(c, http.StatusOK, points)
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

// MarkPDFExported flips the pdf_exported flag for the visitor identified
// by URL param :id. Public route — anyone can mark their own visit, since
// the visitor ID is only known to a session that just received it from
// /api/visitors POST.
func (h *VisitorHandler) MarkPDFExported(c *gin.Context) {
	visitorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.Error(c, http.StatusBadRequest, 40000, "invalid visitor id")
		return
	}

	if err := h.store.MarkPDFExported(visitorID); err != nil {
		httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to mark pdf exported")
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
