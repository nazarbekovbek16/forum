package model

type Item struct {
	ID       int
	Email    string
	Username string
	Password string
}

type Home struct {
	Posts       []PostItem
	CurrentUser bool
}

type Read struct {
	Posts       PostItem
	CurrentUser bool
	Comments    []CommentsItem
}

type PostItem struct {
	ID       int
	Owner    string
	Title    string
	Content  string
	Types    string
	Image    string
	Likes    int
	Dislikes int
}

type TypePostItem struct {
	PostId int
	Type   string
}

type CommentsItem struct {
	ID       int
	PostId   int
	Owner    string
	Comment  string
	Likes    int
	Dislikes int
}
