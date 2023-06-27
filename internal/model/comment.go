package model

type Commentary struct {
	ID       int
	PostID   int
	Content  string
	Author   string
	Likes    int
	Dislikes int
}
