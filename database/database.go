package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

type NewsFeed struct {
	DB *sql.DB
}

var DB *sql.DB

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

func Feed(db *sql.DB) *NewsFeed {
	db, err := sql.Open("sqlite3", "./database/feed.db")
	if err != nil {
		fmt.Printf("Feed db sql.Open error: %+v\n", err)
	}

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS feed (iD INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, title TEXT, content	TEXT, likes INTEGER, created TEXT, category TEXT)")
	stmt.Exec()
	if err != nil {
		fmt.Printf("Feed stmt sql.Open error: %+v\n", err)
	}
	return &NewsFeed{
		DB: db,
	}
}

func (feed *NewsFeed) Get() []PostFeed {
	var id int
	var title string
	var content string
	// var comments []string
	var likes int
	var created string
	var category string

	posts := []PostFeed{}

	rows, err := feed.DB.Query("SELECT * FROM feed")
	if err != nil {
		fmt.Printf("Feed DB Query error: %+v\n", err)
	}

	for rows.Next() {
		rows.Scan(&id, &title, &content, &likes, &created, &category)
		newPost := PostFeed{ // explicit values
			ID:      id,
			Title:   title,
			Content: content,
			Likes:   likes,
			Created: created,
			// Comments: comments,
			Category: category,
		}
		posts = append(posts, newPost)
	}
	return posts
}

func (feed *NewsFeed) Add(item PostFeed) {
	stmt, err := feed.DB.Prepare(`
		INSERT INTO feed (Content) values(?)
	`)
	if err != nil {
		fmt.Printf("feed DB Prepare error: %+v\n", err)
	}
	stmt.Exec(item.Content)
}
