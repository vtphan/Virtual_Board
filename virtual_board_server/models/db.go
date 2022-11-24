package models

import "database/sql"

type DB struct {
	*sql.DB
}

var db *DB

func init() {
	db, _ = InitDB()
}

func InitDB() (*DB, error) {
	db, err := sql.Open("sqlite3", "./Test1.db")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
