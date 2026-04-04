package model

import "time"

type VisitorRecord struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	IP        string    `json:"ip" gorm:"column:ip;not null;index"`
	City      string    `json:"city" gorm:"column:city;not null"`
	Device    string    `json:"device" gorm:"column:device;not null"`
	Browser   string    `json:"browser" gorm:"column:browser;not null"`
	OS        string    `json:"os" gorm:"column:os;not null"`
	VisitTime time.Time `json:"visitTime" gorm:"column:visit_time;not null;index"`
	Duration  int       `json:"duration" gorm:"column:duration;not null;default:0"`
}

func (VisitorRecord) TableName() string {
	return "visitors"
}

type VisitorStats struct {
	TotalVisits        int `json:"totalVisits"`
	TodayVisits        int `json:"todayVisits"`
	UniqueVisitors     int `json:"uniqueVisitors"`
	AverageDurationSec int `json:"averageDurationSeconds"`
}
