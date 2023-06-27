package repository

import (
	"database/sql"
	"forum/internal/config"
)

type Repository struct {
	Auth
	Post
	Commentary
	VotePost
	VoteComment
	User
}

func NewRepository(db *sql.DB, cfg *config.Config) *Repository {
	return &Repository{
		Auth:        newAuthRepository(db, cfg),
		Post:        newPostRepository(db, cfg),
		Commentary:  newCommentaryRepository(db, cfg),
		VotePost:    newVotePostRepository(db, cfg),
		VoteComment: newVoteCommentaryRepository(db, cfg),
		User:        newUserRepository(db, cfg),
	}
}
