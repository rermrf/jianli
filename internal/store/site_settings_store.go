package store

import (
	"gorm.io/gorm"

	"jianli/internal/model"
)

type SiteSettings struct {
	AllowPDFExport bool
}

type SiteSettingsStore struct {
	db *gorm.DB
}

func NewSiteSettingsStore(db *gorm.DB) *SiteSettingsStore {
	return &SiteSettingsStore{db: db}
}

func (s *SiteSettingsStore) Get() (SiteSettings, error) {
	var record model.SiteSettingsRecord
	if err := s.db.First(&record, "id = ?", 1).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return SiteSettings{}, err
		}

		record = model.SiteSettingsRecord{ID: 1, AllowPDFExport: true}
		if err := s.db.Create(&record).Error; err != nil {
			return SiteSettings{}, err
		}
	}

	return SiteSettings{AllowPDFExport: record.AllowPDFExport}, nil
}

func (s *SiteSettingsStore) Save(allowPDFExport bool) error {
	if _, err := s.Get(); err != nil {
		return err
	}

	return s.db.Model(&model.SiteSettingsRecord{}).
		Where("id = ?", 1).
		Update("allow_pdf_export", allowPDFExport).
		Error
}
