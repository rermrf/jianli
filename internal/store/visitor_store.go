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
	err := s.db.Where("ip = ? AND visit_time >= ?", record.IP, record.VisitTime.Add(-10*time.Minute)).Order("visit_time DESC").First(&existing).Error
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
	if days > 0 {
		query = query.Where("visit_time >= ?", time.Now().AddDate(0, 0, -days))
	}

	var records []model.VisitorRecord
	err := query.Order("visit_time DESC").Limit(limit).Offset((page - 1) * limit).Find(&records).Error
	return records, err
}

func (s *VisitorStore) Stats(days int) (model.VisitorStats, error) {
	query := s.db.Model(&model.VisitorRecord{})
	if days > 0 {
		query = query.Where("visit_time >= ?", time.Now().AddDate(0, 0, -days))
	}

	var result struct {
		TotalVisits    int
		UniqueVisitors int
		Average        float64
	}
	if err := query.Select("COUNT(*) AS total_visits, COUNT(DISTINCT ip) AS unique_visitors, COALESCE(AVG(duration), 0) AS average").Scan(&result).Error; err != nil {
		return model.VisitorStats{}, err
	}

	var todayVisits int64
	if err := s.db.Model(&model.VisitorRecord{}).Where("date(visit_time) = date('now', 'localtime')").Count(&todayVisits).Error; err != nil {
		return model.VisitorStats{}, err
	}

	return model.VisitorStats{
		TotalVisits:        result.TotalVisits,
		TodayVisits:        int(todayVisits),
		UniqueVisitors:     result.UniqueVisitors,
		AverageDurationSec: int(result.Average),
	}, nil
}

// Trend returns visit counts bucketed by local-time day for the last N
// days. Days <= 0 is rejected so callers cannot accidentally request a
// full-table scan; the frontend translates 'all' to days=30 explicitly.
//
// Buckets are produced by SQLite using strftime with the 'localtime'
// modifier, so the result depends on the process's TZ. Production sets
// TZ=Asia/Shanghai in the Dockerfile; tests pin time.Local before running.
//
// Sparse days (no visits) are omitted from the result; the frontend pads
// missing dates with count=0 to render a complete N-bar chart.
func (s *VisitorStore) Trend(days int) ([]model.VisitorTrendPoint, error) {
	if days <= 0 {
		return nil, errors.New("trend: days must be positive")
	}

	cutoff := time.Now().AddDate(0, 0, -days+1) // include today's bucket

	type row struct {
		Date  string
		Count int
	}
	var rows []row
	err := s.db.Model(&model.VisitorRecord{}).
		Select("strftime('%Y-%m-%d', visit_time, 'localtime') AS date, COUNT(*) AS count").
		Where("visit_time >= ?", cutoff).
		Group("date").
		Order("date").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("trend query: %w", err)
	}

	points := make([]model.VisitorTrendPoint, len(rows))
	for i, r := range rows {
		points[i] = model.VisitorTrendPoint{Date: r.Date, Count: r.Count}
	}
	return points, nil
}
