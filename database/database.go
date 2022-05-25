package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

type NewsFeed struct {
	DB *sql.DB
}

type CommentFeed struct {
	DB2 *sql.DB
}

var (
	DB  *sql.DB
	DB2 *sql.DB
)

// opens database and checks if table exists, if not makes one
func UserDatabase() {
	db, err := sql.Open("sqlite3", "./database/userdata.db")
	if err != nil {
		fmt.Printf("UserDatabase db sql.Open error: %+v\n", err)
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT, email TEXT, passwordHASH TEXT)")
	statement.Exec()
	if err != nil {
		fmt.Printf("UserDatabase statement sql.Open error: %+v\n", err)
	}
	DB = db
}

// opens and creates a different database
func Feed(db *sql.DB) *NewsFeed {
	db, err := sql.Open("sqlite3", "./database/feed.db")
	if err != nil {
		fmt.Printf("Feed db sql.Open error: %+v\n", err)
	}

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS feed (iD INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, title TEXT, content	TEXT, likes INTEGER, dislikes INTEGER, created TEXT, category TEXT)")
	stmt.Exec()
	if err != nil {
		fmt.Printf("Feed stmt sql.Open error: %+v\n", err)
	}
	return &NewsFeed{
		DB: db,
	}
}

func Comments(db2 *sql.DB) *CommentFeed {
	db, err := sql.Open("sqlite3", "./database/comments.db")
	if err != nil {
		fmt.Printf("Comments Feed sql.Open error: %+v\n", err)
	}

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS comments (iD INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, title TEXT, content	TEXT, likes INTEGER, dislikes INTEGER, created TEXT, category TEXT)")
	stmt.Exec()
	if err != nil {
		fmt.Printf("Comments stmt db.Prepare error: %+v\n", err)
	}
	return &CommentFeed{
		DB2: db,
	}
}

// Get() dumps all values from a selected table
func (feed *NewsFeed) Get() []PostFeed {
	// variable init
	var id int
	var title string
	var content string
	// var comments []string
	var likes int
	var dislikes int
	var created string
	var category string

	posts := []PostFeed{}

	rows, err := feed.DB.Query("SELECT * FROM feed")
	if err != nil {
		fmt.Printf("Feed DB Query error: %+v\n", err)
	}

	// scan rows in database, update variable using memory addresses and link to struct
	for rows.Next() {
		rows.Scan(&id, &title, &content, &likes, &dislikes, &created, &category)
		newPost := PostFeed{ // explicit values
			ID:       id,
			Title:    title,
			Content:  content,
			Likes:    likes,
			Dislikes: dislikes,
			Created:  created,
			// Comments: comments,
			Category: category,
		}
		posts = append(posts, newPost)
	}
	return posts
}

func (feed *CommentFeed) GetComments() []PostFeed {
	// variable init
	var id int
	var title string
	var content string
	// var comments []string
	var likes int
	var dislikes int
	var created string
	var category string

	posts := []PostFeed{}
	rows, err := feed.DB2.Query("SELECT * FROM comments")
	if err != nil {
		fmt.Printf("Comments Feed DB Query error: %+v\n", err)
	}

	// scan rows in database, update variable using memory addresses and link to struct
	for rows.Next() {
		rows.Scan(&id, &title, &content, &likes, &dislikes, &created, &category)
		newComment := PostFeed{ // explicit values
			ID:       id,
			Title:    title,
			Content:  content,
			Likes:    likes,
			Dislikes: dislikes,
			Created:  created,
			// Comments: comments,
			Category: category,
		}
		posts = append(posts, newComment)
	}
	return posts
}

// Add(adds an item into a table)
func (feed *NewsFeed) Add(item PostFeed) {
	stmt, err := feed.DB.Prepare("INSERT INTO feed (title, content, likes, dislikes, created, category) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		fmt.Printf("feed DB Prepare error: %+v\n", err)
	}
	// stmt.QueryRow(stmt, item.Title, item.Content, item.Category)
	stmt.Exec(item.Title, item.Content, item.Likes, item.Dislikes, item.Created, item.Category)
	defer stmt.Close()
}

// Add(adds an item into a table)
func (feed *CommentFeed) AddComment(item PostFeed) {
	stmt, err := feed.DB2.Prepare("INSERT INTO comments (title, content, likes, dislikes, created, category) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		fmt.Printf("AddComment DB Prepare error: %+v\n", err)
	}
	// stmt.QueryRow(stmt, item.Title, item.Content, item.Category)
	stmt.Exec(item.Title, item.Content, item.Likes, item.Dislikes, item.Created, item.Category)
	defer stmt.Close()
}
