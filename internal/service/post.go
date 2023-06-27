package service

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"forum/internal/repository"
	"strings"
)

var (
	ErrInvalidPostTitle   = errors.New("invalid post title characters")
	ErrInvalidPostContent = errors.New("invalid post content characters")
	ErrPostTitleLen       = errors.New("title length out of range")
	ErrPostContentLen     = errors.New("content length out of range")
)

type Post interface {
	CreatePost(post model.Post) error
	GetAllPosts() ([]model.Post, error)
	GetPostByID(postId int) (model.Post, error)
	GetAllPostsByFilter(user model.User, query map[string][]string) ([]model.Post, error)
}

type PostService struct {
	Repository repository.Post
}

func newPostService(repository repository.Post) *PostService {
	return &PostService{
		Repository: repository,
	}
}

func checkPost(post model.Post) error {
	if len(post.Title) > 100 {
		return ErrPostTitleLen
	}

	if len(post.Content) > 1500 {
		return ErrPostContentLen
	}

	post.Title = strings.Trim(post.Title, " \n\r")

	for _, char := range post.Title {
		if (char != 13 && char != 10 && char != 9) && (char < 32 || char > 126) {
			return fmt.Errorf("service: Create Post: check post: %w", ErrInvalidPostTitle)
		}
	}

	if post.Title == "" {
		return fmt.Errorf("service: Create Post: check post: %w", ErrInvalidPostTitle)
	}

	post.Content = strings.Trim(post.Content, " \n\r")
	for _, char := range post.Content {
		if (char != 13 && char != 10 && char != 9) && (char < 32 || char > 126) {
			return fmt.Errorf("service: Create Post: check post: %w", ErrInvalidPostContent)
		}
	}
	if post.Content == "" {
		return fmt.Errorf("service: Create Post: check post: %w", ErrInvalidPostContent)
	}

	return nil
}

func (s *PostService) CreatePost(post model.Post) error {
	if err := checkPost(post); err != nil {
		return err
	}
	if err := s.Repository.CreatePost(post); err != nil {
		return err
	}
	return nil
}

func (s *PostService) GetAllPosts() ([]model.Post, error) {
	allPosts, err := s.Repository.GetAllPosts()
	if err != nil {
		return nil, err
	}
	for i := range allPosts {
		category, err := s.Repository.GetCategoriesByPostID(allPosts[i].ID)
		if err != nil {
			return nil, err
		}
		allPosts[i].Category = category
	}
	return allPosts, nil
}

func (s *PostService) GetPostByID(postId int) (model.Post, error) {
	post, err := s.Repository.GetPostByID(postId)
	if err != nil {
		return model.Post{}, err
	}

	post.Category, err = s.Repository.GetCategoriesByPostID(postId)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (s *PostService) GetAllPostsByFilter(user model.User, query map[string][]string) ([]model.Post, error) {
	var allPosts []model.Post
	var err error

	for key, value := range query {
		switch key {
		case "category":
			allPosts, err = s.Repository.GetPostsByCategory(strings.Join(value, ""))
			if err != nil {
				return nil, err
			}
		case "time":
			switch strings.Join(value, "") {
			case "new":
				allPosts, err = s.Repository.GetNewestPosts()
			case "old":
				allPosts, err = s.Repository.GetOldestPosts()
			default:
				return nil, err
			}
			if err != nil {
				return nil, err
			}
		case "vote":
			switch strings.Join(value, "") {
			case "like":
				allPosts, err = s.Repository.GetMostLikedPosts()
			case "dislike":
				allPosts, err = s.Repository.GetMostDislikedPosts()
			default:
				return nil, err
			}
			if err != nil {
				return nil, err
			}
		case "clean":
			switch strings.Join(value, "") {
			case "true":
				allPosts, err = s.Repository.GetAllPosts()

			}
		default:
			return nil, errors.New("post filter service")
		}

		for i := range allPosts {
			categories, err := s.Repository.GetCategoriesByPostID(allPosts[i].ID)
			if err != nil {
				return nil, err
			}
			allPosts[i].Category = categories
		}
	}
	return allPosts, nil
}
