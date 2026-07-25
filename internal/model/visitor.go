package model

import "time"

type VisitorRecord struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	IP          string    `json:"ip" gorm:"column:ip;not null;index"`
	City        string    `json:"city" gorm:"column:city;not null"`
	Region      string    `json:"region" gorm:"column:region;not null;default:''"`
	Country     string    `json:"country" gorm:"column:country;not null;default:''"`
	ISP         string    `json:"isp" gorm:"column:isp;not null;default:''"`
	Device      string    `json:"device" gorm:"column:device;not null"`
	Browser     string    `json:"browser" gorm:"column:browser;not null"`
	OS          string    `json:"os" gorm:"column:os;not null"`
	VisitTime   time.Time `json:"visitTime" gorm:"column:visit_time;not null;index"`
	Duration    int       `json:"duration" gorm:"column:duration;not null;default:0"`
	PDFExported bool      `json:"pdfExported" gorm:"column:pdf_exported;not null;default:false"`
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

// VisitorTrendPoint is one bucket of the daily visits trend.
type VisitorTrendPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}
