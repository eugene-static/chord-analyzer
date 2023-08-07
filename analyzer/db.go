package analyzer

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3"

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	db, err := sql.Open(driverName, storagePath)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Save(name string, b []byte) error {
	query := "INSERT INTO files VALUES(?, ?)"
	_, err := s.db.Exec(query, name, b)
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) Get(name string) ([]byte, error) {
	query := "SELECT bytes FROM files WHERE name=?"
	var b []byte
	err := s.db.QueryRow(query, name).Scan(&b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
