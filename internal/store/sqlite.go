package store

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	dsn := filepath.Clean(path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := initializeSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func initializeSchema(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS resume (
			id INTEGER PRIMARY KEY DEFAULT 1,
			data TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS visitors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip TEXT NOT NULL,
			city TEXT NOT NULL,
			device TEXT NOT NULL,
			browser TEXT NOT NULL,
			os TEXT NOT NULL,
			visit_time DATETIME NOT NULL,
			duration INTEGER NOT NULL DEFAULT 0
		);`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("initialize schema: %w", err)
		}
	}

	return nil
}
