package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

// type NewsFeed struct {
// 	DB *sql.DB
// }

// type CommentFeed struct {
// 	DB2 *sql.DB
// }

// var (
// 	DB  *sql.DB
// 	DB2 *sql.DB
// )

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


func (forum *Forum) CreateSession(session Session){
	stmt, err := forum.DB.Prepare("INSERT INTO session (uuid, session_uuid) VALUES (?, ?);")
	if err != nil {
		CheckErr(err)
	}
	stmt.Exec(session.Id, session.Uuid)
	defer stmt.Close()



}

func (forum *Forum) CreatePost(post PostFeed) {

	stmt, err := forum.DB.Prepare("INSERT INTO post (postID, authID, title, content, likes, dislikes, category, dateCreated,) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		CheckErr(err)
	}
	stmt.Exec(post.PostID, post.Uuid, post.Title, post.Content, post.Likes, post.Dislikes, post.Category, post.CreatedAt)
	defer stmt.Close()
}



func (forum *Forum) CreateComment(comment Comment){
	stmt, err := forum.DB.Prepare("INSERT INTO comments (commentID, authID, postID, content, dateCreated,) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		CheckErr(err)
	}
	stmt.Exec(comment.CommentID, comment.Uuid, comment.PostID, comment.Content, comment.CreatedAt)
	defer stmt.Close()



}

// ---------------------------------------------- TABLES ---------------------------------//

func userTable(db *sql.DB) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS people (
	uuid TEXT PRIMARY KEY, 
	username TEXT,
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
	authID TEXT,
	session_uuid INTEGER PRIMARY KEY,
	foreign key (authID) references people(uuid));
	`)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()

}

func postTabe(db *sql.DB) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS post (
 postID INTEGER PRIMARY KEY,
 authID TEXT, 
 title TEXT,
 content TEXT, 
 likes INTEGER,
 dislikes INTEGER,
 category TEXT,
 dateCreated TEXT,
 foreign key (authID) references people(uuid));
 `)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()

}

func commentTable(db *sql.DB) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments (
   commentID INTEGER PRIMARY KEY , 
	authID TEXT,
	postID TEXT, 
	content TEXT, 
	dateCreated TEXT, 
	foreign key (authID) references people(uuid),
	foreign key (commentID) references post(postID));
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


// func Comments(db2 *sql.DB) *CommentFeed {
// 	db, err := sql.Open("sqlite3", "./database/comments.db")
// 	if err != nil {
// 		fmt.Printf("Comments Feed sql.Open error: %+v\n", err)
// 	}

// 	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS comments (iD INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, title TEXT, content	TEXT, likes INTEGER, created TEXT, category TEXT)")
// 	stmt.Exec()
// 	if err != nil {
// 		fmt.Printf("Comments stmt db.Prepare error: %+v\n", err)
// 	}
// 	return &CommentFeed{
// 		DB2: db,
// 	}
// }

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

// func (feed *Forum) GetComments() []PostFeed {
// 	// variable init
// 	var id int
// 	var title string
// 	var content string
// 	// var comments []string
// 	var likes int
// 	var created string
// 	var category string

// 	posts := []PostFeed{}

// 	rows, err := feed.DB2.Query("SELECT * FROM comments")
// 	if err != nil {
// 		fmt.Printf("Comments Feed DB Query error: %+v\n", err)
// 	}

	// scan rows in database, update variable using memory addresses and link to struct
// 	for rows.Next() {
// 		rows.Scan(&id, &title, &content, &likes, &created, &category)
// 		newComment := PostFeed{ // explicit values
// 			PostID:    id,
// 			Title:     title,
// 			Content:   content,
// 			Likes:     likes,
// 			CreatedAt: created,
// 			// Comments: comments,
// 			Category: category,
// 		}

// 		posts = append(posts, newComment)
// 	}
// 	return posts
// }

// Add(adds an item into a table)
func (feed *Forum) Add(item PostFeed) {
	stmt, err := feed.DB.Prepare("INSERT INTO feed (title, content, likes, created, category) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		fmt.Printf("feed DB Prepare error: %+v\n", err)
	}
	// stmt.QueryRow(stmt, item.Title, item.Content, item.Category)

	stmt.Exec(item.Title, item.Content, item.Likes, item.CreatedAt, item.Category)

	defer stmt.Close()
}

// func CheckErr(err error) {
// 	fmt.Println(err)
// 	log.Fatal(err)
// }

// Add(adds an item into a table)
func (feed *Forum) AddComment(item PostFeed) {
	stmt, err := feed.DB.Prepare("INSERT INTO comments (title, content, likes, created, category) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		fmt.Printf("AddComment DB Prepare error: %+v\n", err)
	}
	// stmt.QueryRow(stmt, item.Title, item.Content, item.Category)
	stmt.Exec(item.Title, item.Content, item.Likes, item.CreatedAt, item.Category)
	defer stmt.Close()
}
