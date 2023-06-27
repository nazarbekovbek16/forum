package service

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/repository"
)

type VoteComment interface {
	LikeCommentary(commentId int, username string) error
	DislikeCommentary(commentId int, username string) error
	GetCommentaryLikes(postId int) (map[int][]string, error)
	GetCommentaryDislikes(postId int) (map[int][]string, error)
}

type VoteCommentaryService struct {
	Repository repository.VoteComment
}

func newVoteCommentaryService(repository repository.VoteComment) *VoteCommentaryService {
	return &VoteCommentaryService{
		Repository: repository,
	}
}
func (s *VoteCommentaryService) LikeCommentary(commentId int, username string) error {
	if err := s.Repository.CommentaryLiked(commentId, username); err == nil {
		if err := s.Repository.RemoveLikeFromCommentary(commentId, username); err != nil {
			return err
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: like comment: %w", err)
	}

	if err := s.Repository.CommentaryDisliked(commentId, username); err == nil {
		if err := s.Repository.RemoveDislikeFromCommentary(commentId, username); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: like comment: %w", err)
	}

	if err := s.Repository.LikeCommentary(commentId, username); err != nil {
		return err
	}
	return nil
}

func (s *VoteCommentaryService) DislikeCommentary(commentId int, username string) error {
	if err := s.Repository.CommentaryDisliked(commentId, username); err == nil {
		if err := s.Repository.RemoveDislikeFromCommentary(commentId, username); err != nil {
			return err
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: dislike comment: %w", err)
	}

	if err := s.Repository.CommentaryLiked(commentId, username); err == nil {
		if err := s.Repository.RemoveLikeFromCommentary(commentId, username); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: dislike comment: %w", err)
	}

	if err := s.Repository.DislikeCommentary(commentId, username); err != nil {
		return err
	}
	return nil
}

func (s *VoteCommentaryService) GetCommentaryLikes(postId int) (map[int][]string, error) {
	users, err := s.Repository.GetCommentaryLikes(postId)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *VoteCommentaryService) GetCommentaryDislikes(postId int) (map[int][]string, error) {
	users, err := s.Repository.GetCommentaryDislikes(postId)
	if err != nil {
		return nil, err
	}
	return users, nil
}
