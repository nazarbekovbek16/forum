package repository

const (
	TableForUsers = `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		email TEXT UNIQUE,
		username TEXT UNIQUE,
		password TEXT
	);`

	TableForPosts = `CREATE TABLE IF NOT EXISTS post (
		id INTEGER PRIMARY KEY,
		owner TEXT,
		title TEXT,
		content TEXT,
		type TEXT,
		image Text,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0
	);`
	TableForComments = `CREATE TABLE IF NOT EXISTS comment (
		id INTEGER PRIMARY KEY,
		postID INTEGER,
		owner TEXT,
		content TEXT
	);`
	TableForLikes = `Create table if not exists likes(
		id INTEGER PRIMARY KEY,
		postID INTEGER,
		owner TEXT,
		like INTEGER DEFAULT 0,
		dislike INTEGER DEFAULT 0
	);`
	TableForLikesComment = `Create table if not exists likescomment(
		id INTEGER PRIMARY KEY,
		commentID INTEGER,
		owner TEXT,
		like INTEGER DEFAULT 0,
		dislike INTEGER DEFAULT 0
	);`
	TableForSession = `CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY,
		user TEXT,
		session TEXT
	);`
)
