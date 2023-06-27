package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/config"
	"forum/internal/model"
	"time"
)

type Auth interface {
	CreateUser(user model.User) error
	GetUser(username string) (model.User, error)
	SaveToken(username, token string, expirationTime time.Time) error
	GetUserByToken(token string) (model.User, error)
	DeleteToken(token string) error
}
type AuthRepository struct {
	db  *sql.DB
	cfg *config.Config
}

func newAuthRepository(db *sql.DB, cfg *config.Config) *AuthRepository {
	return &AuthRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *AuthRepository) CreateUser(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `INSERT INTO user (email, username, password) VALUES ($1, $2, $3);`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("repository: create user: %w	", err)
	}
	return nil
}

func (r *AuthRepository) GetUser(username string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT id, email, username, password FROM user WHERE username = $1;`
	var user model.User
	if err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		return model.User{}, fmt.Errorf("repository: get user: %w", err)
	}
	return user, nil
}

func (r *AuthRepository) SaveToken(username, token string, expirationTime time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `UPDATE user SET token = $1, expiration_time = $2 WHERE username = $3;`
	_, err := r.db.ExecContext(ctx, query, token, expirationTime, username)
	if err != nil {
		return fmt.Errorf("repository: save token: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetUserByToken(token string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `SELECT id, email, username, password, expiration_time FROM user WHERE token = $1;`
	var user model.User
	if err := r.db.QueryRowContext(ctx, query, token).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.ExpirationTime); err != nil {
		return model.User{}, fmt.Errorf("repository: get user by token: %w", err)
	}
	return user, nil
}

func (r *AuthRepository) DeleteToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Db.CtxTimeout)*time.Second)
	defer cancel()
	query := `UPDATE user SET token = NULL, expiration_time = NULL WHERE token = $1;`
	_, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("repository: delete token: %w", err)
	}
	return nil
}
