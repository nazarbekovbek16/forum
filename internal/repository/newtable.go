package repository

import (
	"database/sql"
)

type Database struct {
	DB *sql.DB
}

func NewTable(db *sql.DB, str string) *Database {
	statement, _ := db.Prepare(str)
	statement.Exec()
	return &Database{
		DB: db,
	}
}
