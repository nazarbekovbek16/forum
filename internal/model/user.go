package model

import "time"

type User struct {
	ID              int
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
	Posts           int

	Token          string
	ExpirationTime time.Time
}
