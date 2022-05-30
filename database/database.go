package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

type Forum struct {
	*sql.DB
}

var DB *sql.DB

func CheckErr(err error) {
	fmt.Println(err)
	log.Fatal(err)
}

func (forum *Forum) CreateUser(user User) {

	stmt, err := forum.DB.Prepare("INSERT INTO people (uuid, username, email, password) VALUES (?, ?, ?, ?);")
	if err != nil {
		CheckErr(err)
	}
	stmt.Exec(user.Uuid, user.Username, user.Email, user.Password)
	defer stmt.Close()
}

func (forum *Forum) CreateSession(session Session) {
	stmt, err := forum.DB.Prepare("INSERT INTO session (expiryTime, userName) VALUES (?, ?);")
	if err != nil {
		CheckErr(err)
	}
	stmt.Exec(session.Expiry, session.Username)
	defer stmt.Close()
}

func (forum *Forum) CreatePost(post PostFeed, user User) {

	stmt, err := forum.DB.Prepare("INSERT INTO post (userID, title, content, likes, dislikes, category, dateCreated) VALUES (?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		CheckErr(err)
	}
	stmt.Exec(user.UserID, post.Title, post.Content, post.Likes, post.Dislikes, post.Category, post.CreatedAt)
	defer stmt.Close()
}

func (forum *Forum) CreateComment(comment Comment) {
	stmt, err := forum.DB.Prepare("INSERT INTO comments (userID, postID, content, dateCreated,) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		CheckErr(err)
	}
	stmt.Exec(comment.Uuid, comment.PostID, comment.Content, comment.CreatedAt)
	defer stmt.Close()

}

// ---------------------------------------------- TABLES ---------------------------------//

func userTable(db *sql.DB) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS people (
	userID INTEGER PRIMARY KEY,	
	uuid TEXT, 
	username TEXT UNIQUE,
	email TEXT UNIQUE, 
	password TEXT);
`)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()
}

func sessionTable(db *sql.DB) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS session (
	userName TEXT PRIMARY KEY REFERENCES people(userID),
	expiryTime TEXT);
	`)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()

}

func postTabe(db *sql.DB) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS post (
 postID INTEGER PRIMARY KEY AUTOINCREMENT,
 userID INTEGER REFERENCES people(userID),
 title TEXT,
 content TEXT, 
 likes INTEGER,
 dislikes INTEGER,
 category TEXT,
 dateCreated TEXT);
 `)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()

}

func commentTable(db *sql.DB) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments (
   commentID TEXT PRIMARY KEY, 
	userID INTEGER REFERENCES people(userID),
	postID INTEGER REFERENCES people(userID), 
	content TEXT NOT NULL, 
	dateCreated TEXT NOT NULL);
	`)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()

}

func Connect(db *sql.DB) *Forum {
	userTable(db)
	sessionTable(db)
	postTabe(db)
	commentTable(db)

	return &Forum{
		DB: db,
	}
}


func (data *Forum) Get() []PostFeed{
	
	posts := []PostFeed{}

	rows, _:= data.DB.Query(`
	SELECT * FROM post
`)

var id int 
var uiD int 
var title string
var content string 
var likes int 
var dislikes int 
var created string 
var category string  

	for rows.Next() {
		rows.Scan(&id,&uiD,&title,&content,&likes,&dislikes,&category,&created)

		posts = append(posts, PostFeed{
		PostID:	id,
		UserID: uiD,
		Title: title,
		Content: content,
		Likes: likes,
		Dislikes: dislikes,
		Category: category,
		CreatedAt: created,

	})
}
	//fmt.Println(posts)
	return posts
}


// Get() dumps all values from a selected table
// func (feed *Forum) Get() []PostFeed {
// 	// variable init
// 	var id int
// 	var title string
// 	var content string
// 	// var comments []string
// 	var likes int
// 	var created string
// 	var category string

// 	posts := []PostFeed{}

// 	rows, err := feed.DB.Query("SELECT * FROM feed")
// 	if err != nil {
// 		fmt.Printf("Feed DB Query error: %+v\n", err)
// 	}

// 	// scan rows in database, update variable using memory addresses and link to struct
// 	for rows.Next() {
// 		rows.Scan(&id, &title, &content, &likes, &created, &category)
// 		newPost := PostFeed{ // explicit values
// 			PostID:    id,
// 			Title:     title,
// 			Content:   content,
// 			Likes:     likes,
// 			CreatedAt: created,
// 			// Comments: comments,
// 			Category: category,
// 		}
// 		posts = append(posts, newPost)
// 	}
// 	return posts
// }

// Add(adds an item into a table)
// func (feed *Forum) Add(item PostFeed) {
// 	stmt, err := feed.DB.Prepare("INSERT INTO feed (title, content, likes, created, category) VALUES (?, ?, ?, ?, ?);")
// 	if err != nil {
// 		fmt.Printf("feed DB Prepare error: %+v\n", err)
// 	}
// 	// stmt.QueryRow(stmt, item.Title, item.Content, item.Category)

// 	stmt.Exec(item.Title, item.Content, item.Likes, item.CreatedAt, item.Category)

// 	defer stmt.Close()
// }
