package database

import "database/sql"

type Log struct {
	Loggedin bool
}

type User struct {
	Username string
	Password string
	Email    string
	UserID   string
}

// could it be used to store data for userprofile and use a single template execution???
// holds details of user session-- used for cookies
type session struct {
	Id     int
	Uuid   string // random value to be stored at the browser
	Email  string
	UserID string
	// CreatedAt	time.Time
}

type usrProfile struct {
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
	Dislikes []Reaction
	Shares   []string
	Userinfo map[string]string
	// custom   string
}

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

type PostFeed struct {
	ID       int    `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Content  string `json:"content,omitempty"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
	Created  string `json:"created,omitempty"`
	Category string `json:"category,omitempty"`
	// Comments []string
}

type NewsFeed struct {
	DB *sql.DB
}

type CommentFeed struct {
	DB2 *sql.DB
}
