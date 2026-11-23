package store

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"jianli/internal/model"
)

type VisitorStore struct {
	db *gorm.DB
}

func NewVisitorStore(db *gorm.DB) *VisitorStore {
	return &VisitorStore{db: db}
}

func (s *VisitorStore) RecordVisit(record model.VisitorRecord) (int64, error) {
	var existing model.VisitorRecord
	// Shared public IPs can belong to different devices or browsers.
	err := s.db.Where("ip = ? AND user_agent = ? AND visit_time >= ?", record.IP, record.UserAgent, record.VisitTime.Add(-10*time.Minute)).Order("visit_time DESC").First(&existing).Error
	if err == nil {
		return existing.ID, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	if err := s.db.Create(&record).Error; err != nil {
		return 0, err
	}

	return record.ID, nil
}

func (s *VisitorStore) UpdateDuration(id int64, duration int) error {
	return s.db.Model(&model.VisitorRecord{}).Where("id = ?", id).Update("duration", duration).Error
}

// UpdateLocation repairs geographic fields without changing visit metadata.
func (s *VisitorStore) UpdateLocation(record model.VisitorRecord) error {
	return s.db.Model(&model.VisitorRecord{}).Where("id = ?", record.ID).Updates(map[string]any{
		"country": record.Country,
		"region":  record.Region,
		"city":    record.City,
		"isp":     record.ISP,
	}).Error
}

// MarkPDFExported flips pdf_exported = true on the visitor row identified
// by id. Idempotent — calling twice is fine. Returns nil even when the row
// is missing (admin already cleared it, or visitorID is stale).
func (s *VisitorStore) MarkPDFExported(id int64) error {
	return s.db.Model(&model.VisitorRecord{}).Where("id = ?", id).Update("pdf_exported", true).Error
}

func (s *VisitorStore) List(days, page, limit int) ([]model.VisitorRecord, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	query := s.db.Model(&model.VisitorRecord{})
	start, end := visitorDayWindow(days, time.Now())
	query = query.Where("julianday(visit_time) < julianday(?)", end.Format(time.RFC3339))
	if days > 0 {
		query = query.Where("julianday(visit_time) >= julianday(?)", start.Format(time.RFC3339))
	}

	var records []model.VisitorRecord
	err := query.Order("visit_time DESC").Limit(limit).Offset((page - 1) * limit).Find(&records).Error
	return records, err
}

func (s *VisitorStore) Stats(days int) (model.VisitorStats, error) {
	return s.statsAt(days, time.Now())
}

func (s *VisitorStore) statsAt(days int, now time.Time) (model.VisitorStats, error) {
	query := s.db.Model(&model.VisitorRecord{})
	start, end := visitorDayWindow(days, now)
	query = query.Where("julianday(visit_time) < julianday(?)", end.Format(time.RFC3339))
	if days > 0 {
		query = query.Where("julianday(visit_time) >= julianday(?)", start.Format(time.RFC3339))
	}

	var result struct {
		TotalVisits    int
		UniqueVisitors int
		Average        float64
	}
	if err := query.Select("COUNT(*) AS total_visits, COUNT(DISTINCT NULLIF(TRIM(ip), '')) AS unique_visitors, COALESCE(AVG(duration), 0) AS average").Scan(&result).Error; err != nil {
		return model.VisitorStats{}, err
	}

	var todayVisits int64
	today, _ := visitorDayWindow(1, now)
	if err := s.db.Model(&model.VisitorRecord{}).
		Where("julianday(visit_time) >= julianday(?) AND julianday(visit_time) < julianday(?)", today.Format(time.RFC3339), end.Format(time.RFC3339)).
		Count(&todayVisits).Error; err != nil {
		return model.VisitorStats{}, err
	}

	return model.VisitorStats{
		TotalVisits:        result.TotalVisits,
		TodayVisits:        int(todayVisits),
		UniqueVisitors:     result.UniqueVisitors,
		AverageDurationSec: int(result.Average),
	}, nil
}

// Analytics use Beijing calendar days, independent of server/browser TZ.
var visitorTimeZone = time.FixedZone("Asia/Shanghai", 8*60*60)

func visitorDayWindow(days int, now time.Time) (start, end time.Time) {
	local := now.In(visitorTimeZone)
	today := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, visitorTimeZone)
	return today.AddDate(0, 0, -days+1), today.AddDate(0, 0, 1)
}

// Trend returns a complete daily series with visit counts and distinct IPs.
// days=0 includes history from the earliest visit through today.
func (s *VisitorStore) Trend(days int) ([]model.VisitorTrendPoint, error) {
	return s.trendAt(days, time.Now())
}

func (s *VisitorStore) trendAt(days int, now time.Time) ([]model.VisitorTrendPoint, error) {
	if days < 0 {
		return nil, errors.New("trend: days must not be negative")
	}
	start, end := visitorDayWindow(days, now)
	query := s.db.Model(&model.VisitorRecord{}).
		Where("julianday(visit_time) < julianday(?)", end.Format(time.RFC3339))
	if days > 0 {
		query = query.Where("julianday(visit_time) >= julianday(?)", start.Format(time.RFC3339))
	}
	var rows []model.VisitorTrendPoint
	err := query.
		Select("date(visit_time, '+8 hours') AS date, COUNT(*) AS count, COUNT(DISTINCT NULLIF(TRIM(ip), '')) AS unique_visitors").
		Group("date(visit_time, '+8 hours')").
		Order("date ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("trend query: %w", err)
	}

	if days == 0 {
		if len(rows) == 0 {
			return []model.VisitorTrendPoint{}, nil
		}
		start, err = time.ParseInLocation("2006-01-02", rows[0].Date, visitorTimeZone)
		if err != nil {
			return nil, fmt.Errorf("trend date: %w", err)
		}
	}
	byDate := make(map[string]model.VisitorTrendPoint, len(rows))
	for _, row := range rows {
		byDate[row.Date] = row
	}
	points := make([]model.VisitorTrendPoint, 0)
	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		date := day.Format("2006-01-02")
		point := byDate[date]
		point.Date = date
		points = append(points, point)
	}
	return points, nil
}
