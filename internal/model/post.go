package model

import "time"

type Post struct {
	ID           int
	Author       string
	Title        string
	Content      string
	CreationTime time.Time
	Category     []string
	Likes        int
	Dislikes     int
}
