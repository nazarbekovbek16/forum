package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/config"
	"forum/internal/model"
	"time"
)

type User interface {
	GetPostByUsername(username string) ([]model.Post, error)
	GetLikedPostByUsername(username string) ([]model.Post, error)
	GetDislikedPostByUsername(username string) ([]model.Post, error)
	GetCommentedPostByUsername(username string) ([]model.Post, error)
	GetAllCategoriesByPostId(postId int) ([]string, error)
	GetUserByUsername(username string) (model.User, error)
}

type UserRepository struct {
	db  *sql.DB
	cfg *config.Config
}

func newUserRepository(db *sql.DB, cfg *config.Config) *UserRepository {
	return &UserRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *UserRepository) GetPostByUsername(username string) ([]model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	var allPosts []model.Post
	query := `SELECT id, author, title, content, creation_time, likes, dislikes FROM post WHERE author = $1;`
	rows, err := r.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("repository: user: get post by username: query - %w", err)
	}
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: user: get post by username: scan - %w", err)
		}
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}

func (r *UserRepository) GetLikedPostByUsername(username string) ([]model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	var allPosts []model.Post
	query := `SELECT id, author, title, content, creation_time, likes, dislikes FROM post WHERE id IN (SELECT postID FROM likes WHERE username = $1);`
	rows, err := r.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("repository: user: get liked post by username: query - %w", err)
	}

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: user: get liked post by username: scan - %w", err)
		}
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}

func (r *UserRepository) GetDislikedPostByUsername(username string) ([]model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	var allPosts []model.Post
	query := `SELECT id, author, title, content, creation_time, likes, dislikes FROM post WHERE id IN (SELECT postID FROM dislikes WHERE username = $1);`
	rows, err := r.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("repository: user: get disliked post by username: query - %w", err)
	}

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: user: get disliked post by username: scan - %w", err)
		}
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}

func (r *UserRepository) GetCommentedPostByUsername(username string) ([]model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	var allPosts []model.Post
	query := `SELECT id, author, title, content, creation_time, likes, dislikes FROM post WHERE id IN (SELECT postID FROM commentary WHERE author = $1);`
	rows, err := r.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("repository: user: get commented post by username: query - %w", err)
	}
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: user: get commented post by username: scan - %w", err)
		}
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}

func (r *UserRepository) GetAllCategoriesByPostId(postId int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	queryCategory := `SELECT category FROM post_category where postID = $1;`
	categoryRows, err := r.db.QueryContext(ctx, queryCategory, postId)
	if err != nil {
		return nil, fmt.Errorf("repository: user: get all categories by post id: query - %w", err)
	}

	var categories []string
	for categoryRows.Next() {
		var category string
		if err := categoryRows.Scan(&category); err != nil {
			return nil, fmt.Errorf("repository: user: get all categories by post id: scan - %w", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *UserRepository) GetUserByUsername(username string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	var user model.User
	query := `SELECT id, email, username, posts FROM user WHERE username = $1;`
	if err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Email, &user.Username, &user.Posts); err != nil {
		return model.User{}, fmt.Errorf("repository: user: get user by username: %w", err)
	}
	return user, nil
}
