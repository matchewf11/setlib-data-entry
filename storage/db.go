package storage

import (
	"database/sql"
	_ "embed"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var dbSchema string

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {

	db, err := sql.Open("sqlite3", filepath.Join("storage", "data.db"))
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(dbSchema)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (str *Storage) Close() {
	str.db.Close()
}

func (str *Storage) InsertProblem(section, diff, prob string) error {
	// fix this
	_, err := str.db.Exec(`
	INSERT INTO problems (section, difficulty, problem, author, topic, question_type) 
	VALUES (?, ?, ?, ?, ?, ?);`, section, diff, prob, "Test Author", "Test Topic", "Test Type")
	if err != nil {
		return err
	}
	return nil
}
