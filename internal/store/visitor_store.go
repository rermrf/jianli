package store

import (
	"database/sql"
	"time"

	"jianli/internal/model"
)

type VisitorStore struct {
	db *sql.DB
}

func NewVisitorStore(db *sql.DB) *VisitorStore {
	return &VisitorStore{db: db}
}

func (s *VisitorStore) RecordVisit(record model.VisitorRecord) (int64, error) {
	var existingID int64
	err := s.db.QueryRow(`
		SELECT id
		FROM visitors
		WHERE ip = ? AND visit_time >= ?
		ORDER BY visit_time DESC
		LIMIT 1
	`, record.IP, record.VisitTime.Add(-10*time.Minute)).Scan(&existingID)
	if err == nil {
		return existingID, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	result, err := s.db.Exec(`
		INSERT INTO visitors (ip, city, device, browser, os, visit_time, duration)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, record.IP, record.City, record.Device, record.Browser, record.OS, record.VisitTime, record.Duration)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (s *VisitorStore) UpdateDuration(id int64, duration int) error {
	_, err := s.db.Exec(`UPDATE visitors SET duration = ? WHERE id = ?`, duration, id)
	return err
}

func (s *VisitorStore) List(days, page, limit int) ([]model.VisitorRecord, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	args := []any{}
	query := `
		SELECT id, ip, city, device, browser, os, visit_time, duration
		FROM visitors
	`
	if days > 0 {
		query += ` WHERE visit_time >= ?`
		args = append(args, time.Now().AddDate(0, 0, -days))
	}
	query += ` ORDER BY visit_time DESC LIMIT ? OFFSET ?`
	args = append(args, limit, (page-1)*limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.VisitorRecord
	for rows.Next() {
		var record model.VisitorRecord
		if err := rows.Scan(&record.ID, &record.IP, &record.City, &record.Device, &record.Browser, &record.OS, &record.VisitTime, &record.Duration); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, rows.Err()
}

func (s *VisitorStore) Stats(days int) (model.VisitorStats, error) {
	args := []any{}
	whereClause := ""
	if days > 0 {
		whereClause = "WHERE visit_time >= ?"
		args = append(args, time.Now().AddDate(0, 0, -days))
	}

	query := `
		SELECT COUNT(*), COUNT(DISTINCT ip), COALESCE(AVG(duration), 0)
		FROM visitors ` + whereClause

	var stats model.VisitorStats
	var avg float64
	if err := s.db.QueryRow(query, args...).Scan(&stats.TotalVisits, &stats.UniqueVisitors, &avg); err != nil {
		return model.VisitorStats{}, err
	}
	stats.AverageDurationSec = int(avg)

	if err := s.db.QueryRow(`
		SELECT COUNT(*)
		FROM visitors
		WHERE date(visit_time) = date('now', 'localtime')
	`).Scan(&stats.TodayVisits); err != nil {
		return model.VisitorStats{}, err
	}

	return stats, nil
}
