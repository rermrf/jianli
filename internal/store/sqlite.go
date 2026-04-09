package store

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"jianli/internal/model"
)

func Open(path string) (*gorm.DB, error) {
	dsn := filepath.Clean(path)
	parentDir := filepath.Dir(dsn)
	if parentDir != "." {
		if err := os.MkdirAll(parentDir, 0o755); err != nil {
			return nil, err
		}
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.ResumeRecord{}, &model.ResumeDraftRecord{}, &model.VisitorRecord{}, &model.SiteSettingsRecord{}); err != nil {
		return nil, err
	}

	return db, nil
}
