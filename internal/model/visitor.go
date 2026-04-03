package model

import "time"

type VisitorRecord struct {
	ID        int64     `json:"id"`
	IP        string    `json:"ip"`
	City      string    `json:"city"`
	Device    string    `json:"device"`
	Browser   string    `json:"browser"`
	OS        string    `json:"os"`
	VisitTime time.Time `json:"visitTime"`
	Duration  int       `json:"duration"`
}

type VisitorStats struct {
	TotalVisits        int `json:"totalVisits"`
	TodayVisits        int `json:"todayVisits"`
	UniqueVisitors     int `json:"uniqueVisitors"`
	AverageDurationSec int `json:"averageDurationSeconds"`
}
