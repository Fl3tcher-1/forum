package database

import "time"

type User struct {
	UserID   int
	Uuid     string
	Username string
	Email    string
	Password string
	//	CreatedAt string
}

type Log struct {
	Loggedin bool
}

type PostFeed struct {
	PostID    int `json:"postid,omitempty"`
	Username  string
	Uuid      string
	Title     string
	Content   string
	Likes     int `json:"likes"`
	Dislikes  int `json:"dislikes"`
	Category  string
	CreatedAt string
}

type Session struct {
	SessionID string
	Username  string
	Expiry    time.Time
	//	UserID    int
	LoggedIn bool
}

type Comment struct {
	CommentID int
	PostID    int
	UserId    int
	Content   string
	CreatedAt string
}

type UsrProfile struct {
	Name string
	// image    *os.Open
	Info     string
	Photo    string
	Gender   string
	Age      int
	Location string
	Posts    []string
	Comments []string
	Likes    []string
	Shares   []string
	Userinfo map[string]string
	// custom   string

}
