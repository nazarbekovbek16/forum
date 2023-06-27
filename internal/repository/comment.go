package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/config"
	"forum/internal/model"
	"time"
)

type Commentary interface {
	CreateCommentary(comment model.Commentary) error
	GetCommentaryByID(id int) (model.Commentary, error)
	GetCommentariesByPostID(postId int) ([]model.Commentary, error)
}

type CommentaryRepository struct {
	db  *sql.DB
	cfg *config.Config
}

func newCommentaryRepository(db *sql.DB, cfg *config.Config) *CommentaryRepository {
	return &CommentaryRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *CommentaryRepository) CreateCommentary(comment model.Commentary) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `INSERT INTO commentary(postID, author, content) VALUES ($1, $2, $3);`
	_, err := r.db.ExecContext(ctx, query, comment.PostID, comment.Author, comment.Content)
	if err != nil {
		return fmt.Errorf("repository: create commentary: Insert query - %w", err)
	}
	return nil
}

func (r *CommentaryRepository) GetCommentaryByID(id int) (model.Commentary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT id, postID, author FROM commentary WHERE id = $1;`
	var commentary model.Commentary
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&commentary.ID, &commentary.PostID, &commentary.Author); err != nil {
		return model.Commentary{}, fmt.Errorf("repository: get commentary: %w", err)
	}
	return commentary, nil
}

func (r *CommentaryRepository) GetCommentariesByPostID(postId int) ([]model.Commentary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT id, postID, author, content, likes, dislikes FROM commentary WHERE postID = $1;`
	rows, err := r.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, fmt.Errorf("repository: get commentaries of the post: query - %w", err)
	}
	var commentaries []model.Commentary
	for rows.Next() {
		var commentary model.Commentary
		if err := rows.Scan(&commentary.ID, &commentary.PostID, &commentary.Author, &commentary.Content, &commentary.Likes, &commentary.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get commentaries of the post: scan - %w", err)
		}
		commentaries = append(commentaries, commentary)
	}
	return commentaries, nil
}
