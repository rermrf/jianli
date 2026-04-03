package store

import (
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
