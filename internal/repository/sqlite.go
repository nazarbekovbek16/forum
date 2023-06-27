package repository

import (
	"database/sql"
	"forum/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

const (
	userTable = `CREATE TABLE IF NOT EXISTS user (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE,
			username TEXT UNIQUE,
			password TEXT,
			posts INT DEFAULT 0,

			token TEXT DEFAULT NULL,
			expiration_time DATETIME DEFAULT NULL
		);`

	postTable = `CREATE TABLE IF NOT EXISTS post (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author TEXT,
			title TEXT,
			content TEXT,
			creation_time DATE DEFAULT (datetime('now','localtime')),

			likes INT DEFAULT 0,
			dislikes INT DEFAULT 0,
			FOREIGN KEY (author) REFERENCES user(username)
		);`

	postCategoryTable = `CREATE TABLE IF NOT EXISTS post_category (
			postID INTEGER,
			category TEXT,
			FOREIGN KEY (postID) REFERENCES post(id) ON DELETE CASCADE
		);`

	commentTable = `CREATE TABLE IF NOT EXISTS commentary (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			postID INTEGER,
			author TEXT,
			content TEXT,
			likes INT DEFAULT 0,
			dislikes INT DEFAULT 0,
			FOREIGN KEY (postID) REFERENCES post(id) ON DELETE CASCADE
		);`

	likesTable = `CREATE TABLE IF NOT EXISTS likes (
			username TEXT,
			postID INTEGER DEFAULT NULL,
			commentaryID INTEGER DEFAULT NULL,
			FOREIGN KEY (postID) REFERENCES post(id) ON DELETE CASCADE,
			FOREIGN KEY (commentaryId) REFERENCES commentary(id) ON DELETE CASCADE
		);`

	dislikeTable = `CREATE TABLE IF NOT EXISTS dislikes (
			username TEXT,
			postID INTEGER DEFAULT NULL,
			commentaryID INTEGER DEFAULT NULL,
			FOREIGN KEY (postID) REFERENCES post(id) ON DELETE CASCADE,
			FOREIGN KEY (commentaryID) REFERENCES commentary(id) ON DELETE CASCADE
		);`
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Db.Driver, cfg.Db.DBName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	allTables := []string{userTable, postTable, postCategoryTable, commentTable, likesTable, dislikeTable}
	for _, eachTable := range allTables {
		_, err := db.Exec(eachTable)
		if err != nil {
			return err
		}
	}
	return nil
}
