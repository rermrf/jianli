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
	metadata := deriveUserAgentMetadata(userAgent)
	ip := c.ClientIP()

	country, region, city, isp, _ := h.geoip.Resolve(ip)

	id, err := h.store.RecordVisit(model.VisitorRecord{
		IP:             ip,
		Country:        country,
		Region:         region,
		City:           city,
		ISP:            isp,
		Device:         metadata.Device,
		Browser:        metadata.Browser,
		OS:             metadata.OS,
		OSVersion:      metadata.OSVersion,
		BrowserVersion: metadata.BrowserVersion,
		UserAgent:      userAgent,
		VisitTime:      visitTime,
		Duration:       request.Duration,
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
	// Fill missing geographic fields while preserving already known values.
	for i := range records {
		record := &records[i]
		if locationKnown(record.Country) && locationKnown(record.Region) && locationKnown(record.City) && locationKnown(record.ISP) {
			continue
		}
		country, region, city, isp, ok := h.geoip.Resolve(record.IP)
		if !ok || (country == "" && region == "" && city == "") {
			continue
		}
		changed := false
		for _, field := range []struct {
			target   *string
			resolved string
		}{
			{&record.Country, country}, {&record.Region, region}, {&record.City, city}, {&record.ISP, isp},
		} {
			if !locationKnown(*field.target) && locationKnown(field.resolved) {
				*field.target = field.resolved
				changed = true
			}
		}
		if !changed {
			continue
		}
		if err := h.store.UpdateLocation(*record); err != nil {
			httpapi.Error(c, http.StatusInternalServerError, 50000, "failed to update visitor location")
			return
		}
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

// Trend returns Beijing-day visits and distinct visitors. days=0 means all.
func (h *VisitorHandler) Trend(c *gin.Context) {
	days := queryInt(c, "days", 7)
	if days < 0 {
		httpapi.Error(c, http.StatusBadRequest, 40000, "days must not be negative")
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

func locationKnown(value string) bool {
	value = strings.TrimSpace(value)
	return value != "" && value != "0" && value != "未知" && !strings.EqualFold(value, "unknown")
}
