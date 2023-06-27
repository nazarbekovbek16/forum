package service

import (
	"forum/internal/repository"
)

type Service struct {
	Auth
	Post
	Commentary
	VotePost
	VoteComment
	User
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Auth:        newAuthService(repository.Auth),
		Post:        newPostService(repository.Post),
		Commentary:  newCommentaryService(repository.Commentary),
		VotePost:    newVotePostService(repository.VotePost),
		VoteComment: newVoteCommentaryService(repository.VoteComment),
		User:        newUserService(repository.User),
	}
}
