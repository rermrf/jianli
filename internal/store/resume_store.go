package store

import (
	"encoding/json"

	"gorm.io/gorm"

	"jianli/internal/model"
	"jianli/internal/seed"
)

type ResumeStore struct {
	db *gorm.DB
}

func NewResumeStore(db *gorm.DB) *ResumeStore {
	return &ResumeStore{db: db}
}

func (s *ResumeStore) Get() (json.RawMessage, error) {
	var record model.ResumeRecord
	err := s.db.First(&record, "id = ?", 1).Error
	if err == gorm.ErrRecordNotFound {
		return seed.DefaultResume(), nil
	}
	if err != nil {
		return nil, err
	}

	return json.RawMessage(record.Data), nil
}

func (s *ResumeStore) Save(data json.RawMessage) error {
	return s.db.Save(&model.ResumeRecord{ID: 1, Data: []byte(data)}).Error
}
