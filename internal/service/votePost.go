package service

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/repository"
)

type VotePost interface {
	LikePost(postId int, username string) error
	DislikePost(postId int, username string) error
	GetPostLikes(postId int) ([]string, error)
	GetPostDislikes(postId int) ([]string, error)
}

type VotePostService struct {
	Repository repository.VotePost
}

func newVotePostService(repository repository.VotePost) *VotePostService {
	return &VotePostService{
		Repository: repository,
	}
}

func (s *VotePostService) LikePost(postId int, username string) error {
	if err := s.Repository.PostLiked(postId, username); err == nil {
		if err := s.Repository.RemoveLikeFromPost(postId, username); err != nil {
			return err
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: like post: %w", err)
	}

	if err := s.Repository.PostDisliked(postId, username); err == nil {
		if err := s.Repository.RemoveDislikeFromPost(postId, username); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: like post: %w", err)
	}

	if err := s.Repository.LikePost(postId, username); err != nil {
		return err
	}
	return nil
}

func (s *VotePostService) DislikePost(postId int, username string) error {
	if err := s.Repository.PostDisliked(postId, username); err == nil {
		if err := s.Repository.RemoveDislikeFromPost(postId, username); err != nil {
			return err
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: dislike post: %w", err)
	}

	if err := s.Repository.PostLiked(postId, username); err == nil {
		if err := s.Repository.RemoveLikeFromPost(postId, username); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: dislike post: %w", err)
	}

	if err := s.Repository.DislikePost(postId, username); err != nil {
		return err
	}
	return nil
}

func (s *VotePostService) GetPostLikes(postId int) ([]string, error) {
	users, err := s.Repository.GetPostLikes(postId)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *VotePostService) GetPostDislikes(postId int) ([]string, error) {
	users, err := s.Repository.GetPostDislikes(postId)
	if err != nil {
		return nil, err
	}
	return users, nil
}
