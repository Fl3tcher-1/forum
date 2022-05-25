package database

type User struct {
	Uuid     string
	Username string
	Email    string
	Password string
	//	CreatedAt string

}

type PostFeed struct {
	PostID    int
	Uuid      string
	Title     string
	Content   string
	Likes     int
	Dislikes  int
	Category  string
	CreatedAt string
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt string
}

type Comment struct {
	CommentID int
	PostID    int
	Uuid      string
	UserId    int
	Content   string
	CreatedAt string
}
