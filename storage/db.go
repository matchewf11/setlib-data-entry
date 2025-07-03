package db

import (
	"database/sql"
	"path/filepath"
)

func New() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", filepath.Join("storage", "data.db"))
	if err != nil {
		return nil, err
	}

	const createTable = `
	CREATE TABLE IF NOT EXISTS problems (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		section TEXT NOT NULL,
		difficulty TEXT NOT NULL,
		problem TEXT NOT NULL
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}

	return db, nil
}
