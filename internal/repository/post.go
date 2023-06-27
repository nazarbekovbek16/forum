package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/config"
	"forum/internal/model"
	"time"
)

type Post interface {
	CreatePost(post model.Post) error
	GetAllPosts() ([]model.Post, error)
	GetPostByID(postId int) (model.Post, error)
	GetPostsByCategory(category string) ([]model.Post, error)
	GetNewestPosts() ([]model.Post, error)
	GetOldestPosts() ([]model.Post, error)
	GetMostLikedPosts() ([]model.Post, error)
	GetMostDislikedPosts() ([]model.Post, error)
	GetCategoriesByPostID(postId int) ([]string, error)
}
type PostRepository struct {
	db  *sql.DB
	cfg *config.Config
}

func newPostRepository(db *sql.DB, cfg *config.Config) *PostRepository {
	return &PostRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *PostRepository) CreatePost(post model.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `INSERT INTO post (author, title, content) VALUES ($1, $2, $3) RETURNING id;`
	var id int
	if err := r.db.QueryRowContext(ctx, query, post.Author, post.Title, post.Content).Scan(&id); err != nil {
		return fmt.Errorf("repository: create post: Insert post query %w", err)
	}
	query = `UPDATE user SET posts = posts + 1 WHERE username = $1;`
	_, err := r.db.ExecContext(ctx, query, post.Author)
	if err != nil {
		return fmt.Errorf("repository: create post: Update post query - %w", err)
	}

	query = `INSERT INTO post_category (postId, category) VALUES ($1, $2);`
	for _, category := range post.Category {
		_, err := r.db.ExecContext(ctx, query, id, category)
		if err != nil {
			return fmt.Errorf("repository: create post: Insert category query - %w", err)
		}
	}
	return nil
}

func (r *PostRepository) GetAllPosts() ([]model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT id, author, title, content, creation_time, likes, dislikes FROM post;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: get all posts: query - %w", err)
	}
	defer rows.Close()

	var allPosts []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get all posts: scan - %w", err)
		}
		allPosts = append(allPosts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allPosts, nil
}

func (r *PostRepository) GetPostByID(postId int) (model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT id, author, title, content, creation_time, likes, dislikes FROM post WHERE id = $1;`
	var post model.Post
	if err := r.db.QueryRowContext(ctx, query, postId).Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
		return model.Post{}, fmt.Errorf("repository: get post by id: %w", err)
	}
	return post, nil
}

func (r *PostRepository) GetPostsByCategory(category string) ([]model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT id, author, title, content, creation_time, likes, dislikes FROM post WHERE id IN (SELECT postId FROM post_category WHERE category = $1);`
	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("repository: get post by category: query - %w", err)
	}
	defer rows.Close()

	var allPosts []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get post by category: scan - %w", err)
		}
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}

func (r *PostRepository) GetNewestPosts() ([]model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT id, author, title, content, creation_time, likes, dislikes FROM post ORDER BY creation_time DESC;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: get newest post: query - %w", err)
	}
	defer rows.Close()

	var allPosts []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get newest post: scan - %w", err)
		}
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}

func (r *PostRepository) GetOldestPosts() ([]model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT id, author, title, content, creation_time, likes, dislikes FROM post ORDER BY creation_time;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: get oldest post: query - %w", err)
	}
	defer rows.Close()

	var allPosts []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get oldest post: scan - %w", err)
		}
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}

func (r *PostRepository) GetMostLikedPosts() ([]model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT id, author, title, content, creation_time, likes, dislikes FROM post ORDER BY likes DESC;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: get liked post: query - %w", err)
	}
	defer rows.Close()

	var allPosts []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get liked post: scan - %w", err)
		}
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}

func (r *PostRepository) GetMostDislikedPosts() ([]model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT id, author, title, content, creation_time, likes, dislikes FROM post ORDER BY dislikes DESC;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: get disliked post: query - %w", err)
	}
	defer rows.Close()

	var allPosts []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get disliked post: scan - %w", err)
		}
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}

func (r *PostRepository) GetCategoriesByPostID(postId int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT category FROM post_category WHERE postID = $1;`
	rows, err := r.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, fmt.Errorf("repository: get categories of the post: query - %w", err)
	}
	defer rows.Close()

	var allCategories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("repository: get categories of the post: scan - %w", err)
		}
		allCategories = append(allCategories, category)
	}
	return allCategories, nil
}
