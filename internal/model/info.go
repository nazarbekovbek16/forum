package model

type Info struct {
	User             User
	ProfileUser      User
	Post             Post
	Posts            []Post
	PostLikes        []string
	PostDislikes     []string
	Commentaries     []Commentary
	CommentsLikes    map[int][]string
	CommentsDislikes map[int][]string
}
