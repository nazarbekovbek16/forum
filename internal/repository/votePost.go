package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/config"
	"time"
)

type VotePost interface {
	LikePost(postId int, username string) error
	DislikePost(postId int, username string) error
	RemoveLikeFromPost(postId int, username string) error
	RemoveDislikeFromPost(postId int, username string) error
	PostLiked(postId int, username string) error
	PostDisliked(postId int, username string) error
	GetPostLikes(postId int) ([]string, error)
	GetPostDislikes(postId int) ([]string, error)
}

type VotePostRepository struct {
	db  *sql.DB
	cfg *config.Config
}

func newVotePostRepository(db *sql.DB, cfg *config.Config) *VotePostRepository {
	return &VotePostRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *VotePostRepository) LikePost(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `INSERT INTO likes(username, postID) VALUES ($1, $2);`
	_, err := r.db.ExecContext(ctx, query, username, postId)
	if err != nil {
		return fmt.Errorf("repository: like post: Insert query - %w", err)
	}

	query = `UPDATE post SET likes = likes + 1 WHERE id = $1;`
	_, err = r.db.ExecContext(ctx, query, postId)
	if err != nil {
		return fmt.Errorf("repository: like post: Update query - %w", err)
	}
	return nil
}

func (r *VotePostRepository) DislikePost(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `INSERT INTO dislikes(username, postID) VALUES ($1, $2);`
	_, err := r.db.ExecContext(ctx, query, username, postId)
	if err != nil {
		return fmt.Errorf("repository: dislike post: Insert query - %w", err)
	}

	query = `UPDATE post SET dislikes = dislikes + 1 WHERE id = $1;`
	_, err = r.db.ExecContext(ctx, query, postId)
	if err != nil {
		return fmt.Errorf("repository: dislike post: Update query - %w", err)
	}
	return nil
}

func (r *VotePostRepository) RemoveLikeFromPost(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `DELETE FROM likes WHERE postID = $1 AND username = $2;`
	_, err := r.db.ExecContext(ctx, query, postId, username)
	if err != nil {
		return fmt.Errorf("repository: remove like from post: Delete query - %w", err)
	}

	query = `UPDATE post SET likes = likes - 1 WHERE id = $1;`
	_, err = r.db.ExecContext(ctx, query, postId)
	if err != nil {
		return fmt.Errorf("repository: remove like from post: Update query - %w", err)
	}
	return nil
}

func (r *VotePostRepository) RemoveDislikeFromPost(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `DELETE FROM dislikes WHERE postID = $1 AND username = $2;`
	_, err := r.db.ExecContext(ctx, query, postId, username)
	if err != nil {
		return fmt.Errorf("repository: remove dislike from post: Delete query - %w", err)
	}

	query = `UPDATE post SET dislikes = dislikes - 1 WHERE id = $1;`
	_, err = r.db.ExecContext(ctx, query, postId)
	if err != nil {
		return fmt.Errorf("repository: remove dislike from post: Update query - %w", err)
	}
	return nil
}

func (r *VotePostRepository) PostLiked(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	var u, query string
	query = `SELECT username FROM likes WHERE postID = $1 AND username = $2;`
	if err := r.db.QueryRowContext(ctx, query, postId, username).Scan(&u); err != nil {
		return fmt.Errorf("repository: post liked: %w", err)
	}
	return nil
}

func (r *VotePostRepository) PostDisliked(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	var u, query string
	query = `SELECT username FROM dislikes WHERE postID = $1 AND username = $2;`
	if err := r.db.QueryRowContext(ctx, query, postId, username).Scan(&u); err != nil {
		return fmt.Errorf("repository: post disliked: %w", err)
	}
	return nil
}

func (r *VotePostRepository) GetPostLikes(postId int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	var postLikes []string
	query := `SELECT username FROM likes WHERE postID = $1;`
	rows, err := r.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, fmt.Errorf("repository: get post likes: query - %w", err)
	}
	for rows.Next() {
		var postLike string
		if err := rows.Scan(&postLike); err != nil {
			return nil, fmt.Errorf("repository: get post likes: scan - %w", err)
		}
		postLikes = append(postLikes, postLike)
	}
	return postLikes, nil
}

func (r *VotePostRepository) GetPostDislikes(postId int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	var postDislikes []string
	query := `SELECT username FROM dislikes WHERE postID = $1;`
	rows, err := r.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, fmt.Errorf("repository: get post dislikes: query - %w", err)
	}
	for rows.Next() {
		var postDislike string
		if err := rows.Scan(&postDislike); err != nil {
			return nil, fmt.Errorf("repository: get post dislikes: scan - %w", err)
		}
		postDislikes = append(postDislikes, postDislike)
	}
	return postDislikes, nil
}
