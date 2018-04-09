package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	DB *sql.DB
}

func NewSqlite(path string) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", path)
	sqlite3 := &Sqlite{DB: db}
	return sqlite3, err
}
