package service

import (
	"errors"
	"forum/internal/model"
	"forum/internal/repository"
	"strings"
)

var ErrInvalidQuery = errors.New("invalid query request")

type User interface {
	GetPostByUsername(username string, query map[string][]string) ([]model.Post, error)
	GetUserByUsername(username string) (model.User, error)
}

type UserService struct {
	Repository repository.User
}

func newUserService(repository repository.User) *UserService {
	return &UserService{
		Repository: repository,
	}
}

func (s *UserService) GetPostByUsername(username string, query map[string][]string) ([]model.Post, error) {
	var posts []model.Post
	var err error

	search, ok := query["posts"]
	if !ok {
		return nil, ErrInvalidQuery
	}

	switch strings.Join(search, "") {
	case "created":
		posts, err = s.Repository.GetPostByUsername(username)
	case "liked":
		posts, err = s.Repository.GetLikedPostByUsername(username)
	case "disliked":
		posts, err = s.Repository.GetDislikedPostByUsername(username)
	case "commented":
		posts, err = s.Repository.GetCommentedPostByUsername(username)
	default:
		return nil, ErrInvalidQuery
	}
	if err != nil {
		return nil, err
	}

	for i := range posts {
		category, err := s.Repository.GetAllCategoriesByPostId(posts[i].ID)
		if err != nil {
			return nil, err
		}
		posts[i].Category = category
	}
	return posts, nil
}

func (s *UserService) GetUserByUsername(username string) (model.User, error) {
	return s.Repository.GetUserByUsername(username)
}
