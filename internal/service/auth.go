package service

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/model"
	"forum/internal/repository"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidUsernameLen  = errors.New("username length out of range 32")
	ErrInvalidUsernameChar = errors.New("invalid username characters")
	ErrConfirmPassword     = errors.New("password doesn't match")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserExist           = errors.New("user already exists")
)

type Auth interface {
	CreateUser(user model.User) error
	GenerateToken(username, password string) (model.User, error)
	ParseToken(token string) (model.User, error)
	DeleteToken(token string) error
}
type AuthService struct {
	Repository repository.Auth
}

func newAuthService(repository repository.Auth) *AuthService {
	return &AuthService{
		Repository: repository,
	}
}

func generateHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("service: generate hashed password: %w", err)
	}
	return string(hash), nil
}

func compareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func checkUser(user model.User) error {
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return fmt.Errorf("service: CreateUser: checkUser err: %w", ErrInvalidEmail)
	}

	for _, char := range user.Username {
		if char < 32 || char > 126 {
			return fmt.Errorf("service: CreateUser: checkUser err: %w", ErrInvalidUsernameChar)
		}
	}

	if len(user.Username) < 1 || len(user.Username) >= 36 {
		return fmt.Errorf("service: CreateUser: checkUser err: %w", ErrInvalidUsernameLen)
	}

	if user.Password != user.ConfirmPassword {
		return fmt.Errorf("service: CreateUser: checkUser err: %w", ErrConfirmPassword)
	}
	return nil
}

func (s *AuthService) CreateUser(user model.User) error {
	if _, err := s.Repository.GetUser(user.Username); err == nil {
		return fmt.Errorf("service: CreateUser: get user: %w", ErrUserExist)
	}

	if err := checkUser(user); err != nil {
		return err
	}

	var err error
	user.Password, err = generateHashPassword(user.Password)
	if err != nil {
		return err
	}

	return s.Repository.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (model.User, error) {
	user, err := s.Repository.GetUser(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}

	if err := compareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.User{}, fmt.Errorf("service: compare hash and password: %w: %w", err, ErrUserNotFound)
	}
	user.Token = uuid.NewString()
	user.ExpirationTime = time.Now().Add(12 * time.Hour)

	if err := s.Repository.SaveToken(user.Username, user.Token, user.ExpirationTime); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *AuthService) ParseToken(token string) (model.User, error) {
	user, err := s.Repository.GetUserByToken(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}
	return user, nil
}

func (s *AuthService) DeleteToken(token string) error {
	return s.Repository.DeleteToken(token)
}
