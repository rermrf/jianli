package model

import "time"

type ResumePayload []byte

type ResumeRecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement:false"`
	Data      []byte    `gorm:"column:data;type:TEXT;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ResumeRecord) TableName() string {
	return "resume"
}
