package store

import (
	"encoding/json"

	"gorm.io/gorm"

	"jianli/internal/model"
)

type ResumeDraftStore struct {
	db          *gorm.DB
	resumeStore *ResumeStore
}

func NewResumeDraftStore(db *gorm.DB, resumeStore *ResumeStore) *ResumeDraftStore {
	return &ResumeDraftStore{db: db, resumeStore: resumeStore}
}

func (s *ResumeDraftStore) Create(name string, note string, data json.RawMessage) (int64, error) {
	record := model.ResumeDraftRecord{
		Name: name,
		Note: note,
		Data: []byte(data),
	}
	if err := s.db.Create(&record).Error; err != nil {
		return 0, err
	}

	return record.ID, nil
}

func (s *ResumeDraftStore) List() ([]model.ResumeDraftRecord, error) {
	var records []model.ResumeDraftRecord
	if err := s.db.Order("updated_at DESC").Order("id DESC").Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (s *ResumeDraftStore) Get(id int64) (model.ResumeDraftRecord, error) {
	var record model.ResumeDraftRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return model.ResumeDraftRecord{}, err
	}

	return record, nil
}

func (s *ResumeDraftStore) Publish(id int64) error {
	record, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.resumeStore.Save(json.RawMessage(record.Data))
}

func (s *ResumeDraftStore) Delete(id int64) error {
	result := s.db.Delete(&model.ResumeDraftRecord{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
