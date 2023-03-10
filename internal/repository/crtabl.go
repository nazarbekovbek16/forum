package repository

import "database/sql"

func CreateTable() {
	database, _ := sql.Open("sqlite3", "./forum.db")
	NewTable(database, TableForComments)
	NewTable(database, TableForPosts)
	NewTable(database, TableForUsers)
	NewTable(database, TableForSession)
	NewTable(database, TableForLikes)
	NewTable(database, TableForLikesComment)
}
