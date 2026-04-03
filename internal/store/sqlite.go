package store

import (
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"jianli/internal/model"
)

func Open(path string) (*gorm.DB, error) {
	dsn := filepath.Clean(path)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.ResumeRecord{}, &model.VisitorRecord{}); err != nil {
		return nil, err
	}

	return db, nil
}
