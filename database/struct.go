package database

import (
	"database/sql"
	"time"
)

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
	Likes    []Reaction
	Shares   []string
	Userinfo map[string]string
	// custom   string
}

// holds details of user session-- used for cookies.
type Post struct {
	Title    string
	Content  string
	Date     string
	Comments int
	PostID   string
	UserID   string
	Reaction Reaction
}

type Reaction struct {
	PostID     string
	UserID     string
	ReactionID string
	CommentID  string
	// React      int
	Likes    int
	Dislikes int
}

type Forum struct {
	*sql.DB
}

type CategoryPost struct { // create a []post in order to store multiple posts
	Post []PostFeed
}

// Databases holds our post and comment databases
type Databases struct {
	Post    PostFeed
	Comment []Comment
}
