package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

type NewsFeed struct {
	DB *sql.DB
}

var DB *sql.DB

func UserDatabase() {
	db, err := sql.Open("sqlite3", "./database/userdata.db")
	if err != nil {
		CheckErr(err)
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT, email TEXT, passwordHASH TEXT)")
	statement.Exec()

	if err != nil {
		CheckErr(err)
	}

	DB = db
}

func Feed(db *sql.DB) *NewsFeed {
	db, err := sql.Open("sqlite3", "./database/feed.db")
	if err != nil {
		CheckErr(err)
	}

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS feed (iD INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, title TEXT, content	TEXT, likes INTEGER, created TEXT, category TEXT)")
	stmt.Exec()
	if err != nil {
		CheckErr(err)
	}
	return &NewsFeed{
		DB: db,
	}
}

func (feed *NewsFeed) Get() []PostFeed {

	posts := []PostFeed{}
	
	rows, err := feed.DB.Query("SELECT * FROM feed")

	var id int
	var title string
	var content string
	// var comments []string
	var likes int
	var created string
	var category string

	for rows.Next() {
		rows.Scan(&id, &title, &content, &likes, &created, &category)
		newPost := PostFeed{ //explicit values
			ID:      id,
			Title:   title,
			Content: content,
			Likes:    likes,
			Created: created,
			// Comments: comments,
			Category: category,
		}

		posts = append(posts, newPost)
	}

	if err != nil {
		CheckErr(err)
	}

	return posts
}

func (feed *NewsFeed) Add(item PostFeed) {
	stmt, err := feed.DB.Prepare("INSERT INTO feed (title, content, likes, created, category) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		CheckErr(err)
	}
	stmt.QueryRow(stmt, item.Title, item.Content, item.Category)

	stmt.Exec(item.Title, item.Content, item.Likes, item.Created, item.Category)

	defer stmt.Close()
}

func CheckErr(err error) {
	fmt.Println(err)
	log.Fatal(err)
}
