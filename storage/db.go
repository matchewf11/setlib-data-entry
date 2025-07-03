package db

import (
	"database/sql"
	_ "embed"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var dbSchema string

func New() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", filepath.Join("storage", "data.db"))
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(dbSchema)
	if err != nil {
		return nil, err
	}

	return db, nil
}
