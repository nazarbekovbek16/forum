package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/config"
	"time"
)

type VoteComment interface {
	LikeCommentary(commentaryId int, username string) error
	DislikeCommentary(commentaryId int, username string) error
	RemoveLikeFromCommentary(commentaryId int, username string) error
	RemoveDislikeFromCommentary(commentaryId int, username string) error
	CommentaryLiked(commentaryId int, username string) error
	CommentaryDisliked(commentaryId int, username string) error
	GetCommentaryLikes(postId int) (map[int][]string, error)
	GetCommentaryDislikes(postId int) (map[int][]string, error)
}
type VoteCommentaryRepository struct {
	db  *sql.DB
	cfg *config.Config
}

func newVoteCommentaryRepository(db *sql.DB, cfg *config.Config) *VoteCommentaryRepository {
	return &VoteCommentaryRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *VoteCommentaryRepository) LikeCommentary(commentaryId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `INSERT INTO likes(username, commentaryID) VALUES ($1, $2);`
	_, err := r.db.ExecContext(ctx, query, username, commentaryId)
	if err != nil {
		return fmt.Errorf("repository: like commentary: Insert query - %w", err)
	}

	query = `UPDATE commentary SET likes = likes + 1 WHERE id = $1;`
	_, err = r.db.ExecContext(ctx, query, commentaryId)
	if err != nil {
		return fmt.Errorf("repository: like commentary: Update query - %w", err)
	}
	return nil
}

func (r *VoteCommentaryRepository) DislikeCommentary(commentaryId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `INSERT INTO dislikes(username, commentaryID) VALUES ($1, $2);`
	_, err := r.db.ExecContext(ctx, query, username, commentaryId)
	if err != nil {
		return fmt.Errorf("repository: like commentary: Insert query - %w", err)
	}

	query = `UPDATE commentary SET dislikes = dislikes + 1 WHERE id = $1;`
	_, err = r.db.ExecContext(ctx, query, commentaryId)
	if err != nil {
		return fmt.Errorf("repository: like commentary: Update query - %w", err)
	}
	return nil
}

func (r *VoteCommentaryRepository) RemoveLikeFromCommentary(commentaryId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `DELETE FROM likes WHERE commentaryID = $1 AND username = $2;`
	_, err := r.db.ExecContext(ctx, query, commentaryId, username)
	if err != nil {
		return fmt.Errorf("repository: remove like from commentary: Delete query - %w", err)
	}

	query = `UPDATE commentary SET likes = likes - 1 WHERE id = $1;`
	_, err = r.db.ExecContext(ctx, query, commentaryId)
	if err != nil {
		return fmt.Errorf("repository: remove like from commentary: Update query - %w", err)
	}
	return nil
}

func (r *VoteCommentaryRepository) RemoveDislikeFromCommentary(commentaryId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `DELETE FROM dislikes WHERE commentaryID = $1 AND username = $2;`
	_, err := r.db.ExecContext(ctx, query, commentaryId, username)
	if err != nil {
		return fmt.Errorf("repository: remove dislike from commentary: Delete query - %w", err)
	}

	query = `UPDATE commentary SET dislikes = dislikes - 1 WHERE id = $1;`
	_, err = r.db.ExecContext(ctx, query, commentaryId)
	if err != nil {
		return fmt.Errorf("repository: remove dislike from commentary: Update query - %w", err)
	}
	return nil
}

func (r *VoteCommentaryRepository) CommentaryLiked(commentaryId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	var u, query string
	query = `SELECT username FROM likes WHERE commentaryID = $1 AND username = $2;`
	if err := r.db.QueryRowContext(ctx, query, commentaryId, username).Scan(&u); err != nil {
		return fmt.Errorf("repository: is comment liked: %w", err)
	}
	return nil
}

func (r *VoteCommentaryRepository) CommentaryDisliked(commentaryId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	var u, query string
	query = `SELECT username FROM dislikes WHERE commentaryID = $1 AND username = $2;`
	if err := r.db.QueryRowContext(ctx, query, commentaryId, username).Scan(&u); err != nil {
		return fmt.Errorf("repository: is comment disliked: %w", err)
	}
	return nil
}

func (r *VoteCommentaryRepository) GetCommentaryLikes(postId int) (map[int][]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	queryForCommentsId := `SELECT id FROM commentary WHERE postID = $1;`
	queryForUsers := `SELECT username FROM likes WHERE commentaryID = $1;`
	users := make(map[int][]string)
	rowsComment, err := r.db.QueryContext(ctx, queryForCommentsId, postId)
	if err != nil {
		return nil, fmt.Errorf("repository: get commentary likes: query for comment - %w", err)
	}
	for rowsComment.Next() {
		var id int
		if err := rowsComment.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("repository: get commentary likes: scan for comment - %w", err)
		}
		var usernames []string
		rowsUsers, err := r.db.QueryContext(ctx, queryForUsers, id)
		if err != nil {
			return nil, fmt.Errorf("repository: get commentary likes: query for users - %w", err)
		}
		for rowsUsers.Next() {
			var username string
			if err := rowsUsers.Scan(&username); err != nil {
				return nil, fmt.Errorf("repository: get commentary likes: scan for users - %w", err)
			}
			usernames = append(usernames, username)
		}
		users[id] = usernames
	}
	return users, nil
}

func (r *VoteCommentaryRepository) GetCommentaryDislikes(postId int) (map[int][]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	queryForCommentsId := `SELECT id FROM commentary WHERE postID = $1;`
	queryForUsers := `SELECT username FROM dislikes WHERE commentaryID = $1;`
	users := make(map[int][]string)
	rowsComment, err := r.db.QueryContext(ctx, queryForCommentsId, postId)
	if err != nil {
		return nil, fmt.Errorf("repository: get commentary dislikes: query for comment - %w", err)
	}
	for rowsComment.Next() {
		var id int
		if err := rowsComment.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("repository: get commentary dislikes: scan for comment - %w", err)
		}
		var usernames []string
		rowsUsers, err := r.db.QueryContext(ctx, queryForUsers, id)
		if err != nil {
			return nil, fmt.Errorf("repository: get commentary dislikes: query for users - %w", err)
		}
		for rowsUsers.Next() {
			var username string
			if err := rowsUsers.Scan(&username); err != nil {
				return nil, fmt.Errorf("repository: get commentary dislikes: scan for users - %w", err)
			}
			usernames = append(usernames, username)
		}
		users[id] = usernames
	}
	return users, nil
}
