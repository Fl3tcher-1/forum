package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

type NewsFeed struct {
	DB *sql.DB
}

type CommentFeed struct{
	DB2 *sql.DB
}

var DB *sql.DB
var DB2 *sql.DB

//opens database and checks if table exists, if not makes one
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
//opens and creates a different database
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
func Comments(db2 *sql.DB) *CommentFeed {
	db, err := sql.Open("sqlite3", "./database/comments.db")
	if err != nil {
		CheckErr(err)
	}

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS comments (iD INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, title TEXT, content	TEXT, likes INTEGER, created TEXT, category TEXT)")
	stmt.Exec()
	if err != nil {
		CheckErr(err)
	}
	return &CommentFeed{
		DB2: db,
	}
}

//Get() dumps all values from a selected table 
func (feed *NewsFeed) Get() []PostFeed {
<<<<<<< HEAD
=======

	posts := []PostFeed{}
	
	rows, err := feed.DB.Query("SELECT * FROM feed")

	//variable init
>>>>>>> master
	var id int
	var title string
	var content string
	// var comments []string
	var likes int
	var created string
	var category string

<<<<<<< HEAD
	posts := []PostFeed{}

	rows, err := feed.DB.Query("SELECT * FROM feed")
	if err != nil {
		fmt.Printf("Feed DB Query error: %+v\n", err)
	}

=======
	//scan rows in database, update variable using memory addresses and link to struct 
>>>>>>> master
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

<<<<<<< HEAD
func (feed *NewsFeed) Add(item PostFeed) {
	stmt, err := feed.DB.Prepare(`
		INSERT INTO feed (Content) values(?)
	`)
=======
func (feed *CommentFeed) GetComments() []PostFeed {

	posts := []PostFeed{}
	
	rows, err := feed.DB2.Query("SELECT * FROM comments")

	//variable init
	var id int
	var title string
	var content string
	// var comments []string
	var likes int
	var created string
	var category string

	//scan rows in database, update variable using memory addresses and link to struct 
	for rows.Next() {
		rows.Scan(&id, &title, &content, &likes, &created, &category)
		newComment := PostFeed{ //explicit values
			ID:      id,
			Title:   title,
			Content: content,
			Likes:    likes,
			Created: created,
			// Comments: comments,
			Category: category,
		}

		posts = append(posts, newComment)
	}

>>>>>>> master
	if err != nil {
		fmt.Printf("feed DB Prepare error: %+v\n", err)
	}
<<<<<<< HEAD
	stmt.Exec(item.Content)
}
=======

	return posts
}

//Add(adds an item into a table)
func (feed *NewsFeed) Add(item PostFeed) {
	stmt, err := feed.DB.Prepare("INSERT INTO feed (title, content, likes, created, category) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		CheckErr(err)
	}
	// stmt.QueryRow(stmt, item.Title, item.Content, item.Category)

	stmt.Exec(item.Title, item.Content, item.Likes, item.Created, item.Category)

	defer stmt.Close()
}

func CheckErr(err error) {
	fmt.Println(err)
	log.Fatal(err)
}

//Add(adds an item into a table)
func (feed *CommentFeed) AddComment(item PostFeed) {
	stmt, err := feed.DB2.Prepare("INSERT INTO comments (title, content, likes, created, category) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		CheckErr(err)
	}
	// stmt.QueryRow(stmt, item.Title, item.Content, item.Category)

	stmt.Exec(item.Title, item.Content, item.Likes, item.Created, item.Category)

	defer stmt.Close()
}

>>>>>>> master
