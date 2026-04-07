package model

import "time"

type ResumeDraftRecord struct {
	ID        int64     `gorm:"primaryKey"`
	Name      string    `gorm:"column:name;type:TEXT;not null"`
	Note      string    `gorm:"column:note;type:TEXT;not null;default:''"`
	Data      []byte    `gorm:"column:data;type:TEXT;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ResumeDraftRecord) TableName() string {
	return "resume_drafts"
}
