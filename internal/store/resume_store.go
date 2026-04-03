package store

import (
	"database/sql"
	"encoding/json"

	"jianli/internal/seed"
)

type ResumeStore struct {
	db *sql.DB
}

func NewResumeStore(db *sql.DB) *ResumeStore {
	return &ResumeStore{db: db}
}

func (s *ResumeStore) Get() (json.RawMessage, error) {
	var data []byte
	err := s.db.QueryRow(`SELECT data FROM resume WHERE id = 1`).Scan(&data)
	if err == sql.ErrNoRows {
		seedData := seed.DefaultResume()
		if err := s.Save(seedData); err != nil {
			return nil, err
		}
		return seedData, nil
	}
	if err != nil {
		return nil, err
	}

	return json.RawMessage(data), nil
}

func (s *ResumeStore) Save(data json.RawMessage) error {
	_, err := s.db.Exec(`
		INSERT INTO resume (id, data, updated_at)
		VALUES (1, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(id) DO UPDATE SET data = excluded.data, updated_at = CURRENT_TIMESTAMP
	`, []byte(data))
	return err
}
