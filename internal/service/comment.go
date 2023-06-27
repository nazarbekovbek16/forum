package service

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"forum/internal/repository"
	"strings"
)

var (
	ErrInvalidComment     = errors.New("invalid comment")
	ErrInvalidCommentChar = errors.New("invalid characters")
	ErrCommentLen         = errors.New("comment length out of range")
)

type Commentary interface {
	CreateCommentary(comment model.Commentary) error
	GetCommentaryById(commentId int) (model.Commentary, error)
	GetCommentariesByPostID(postId int) ([]model.Commentary, error)
}

type CommentaryService struct {
	Repository repository.Commentary
}

func newCommentaryService(repository repository.Commentary) *CommentaryService {
	return &CommentaryService{
		Repository: repository,
	}
}

func checkCommentary(comment model.Commentary) error {
	if len(comment.Content) > 700 {
		return fmt.Errorf("service: create comment: %w", ErrCommentLen)
	}

	comment.Content = strings.Trim(comment.Content, " \n\r")

	for _, char := range comment.Content {
		if (char != 13 && char != 10) && (char < 32 || char > 126) {
			return fmt.Errorf("service: Create Comment: check comment err: %w", ErrInvalidCommentChar)
		}
	}

	if comment.Content == "" {
		return fmt.Errorf("service: Create Comment: check comment err: %w", ErrInvalidCommentChar)
	}

	return nil
}

func (s *CommentaryService) CreateCommentary(comment model.Commentary) error {
	if err := checkCommentary(comment); err != nil {
		return err
	}

	if err := s.Repository.CreateCommentary(comment); err != nil {
		return err
	}
	return nil
}

func (s *CommentaryService) GetCommentaryById(commentId int) (model.Commentary, error) {
	commentary, err := s.Repository.GetCommentaryByID(commentId)
	if err != nil {
		return model.Commentary{}, err
	}
	return commentary, nil
}

func (s *CommentaryService) GetCommentariesByPostID(postId int) ([]model.Commentary, error) {
	commentaries, err := s.Repository.GetCommentariesByPostID(postId)
	if err != nil {
		return nil, err
	}
	return commentaries, nil
}
