package database

import (
	"database/sql"
	"time"
	"github.com/dgrijalva/jwt-go"
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
	UserId    string
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

type Forum struct {
	*sql.DB
}

// holds details of user session-- used for cookies.
type Post struct {
	Title    string
	Content  string
	Date     string
	Comments int
}

type CategoryPost struct { // create a []post in order to store multiple posts
	Post []PostFeed
}

//Databases holds our post and comment databases
type Databases struct {
	Post    PostFeed
	Comment []Comment
}

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name`
	LastName      string `json:"family_name`
	jwt.StandardClaims
}

