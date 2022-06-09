package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	stmt, err := forum.DB.Prepare("INSERT INTO session (sessionID, username, expiryTime) VALUES (?, ?, ?);")
	if err != nil {
		CheckErr(err)
	}
	stmt.Exec(session.SessionID, session.Username, session.Expiry)
	defer stmt.Close()
}

func (forum *Forum) CreatePost(post PostFeed) {

	stmt, err := forum.DB.Prepare("INSERT INTO post (username, title, content, likes, dislikes, category, dateCreated) VALUES (?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		CheckErr(err)
	}
	stmt.Exec(post.Username, post.Title, post.Content, post.Likes, post.Dislikes, post.Category, post.CreatedAt)
	defer stmt.Close()
}

func (forum *Forum) CreateComment(comment Comment) {
	stmt, err := forum.DB.Prepare("INSERT INTO comments ( postID, userID, content, dateCreated) VALUES (?, ?, ?, ?);")
	if err != nil {
		CheckErr(err)
	}
	stmt.Exec(comment.PostID, comment.UserId, comment.Content,comment.CreatedAt)
	defer stmt.Close()

}

// ---------------------------------------------- TABLES ------------------------------- --//

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
	sessionID TEXT PRIMARY KEY REFERENCES people(uuid),	
	username TEXT REFERENCES people(username), 
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
 username TEXT REFERENCES people(username),
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
   commentID INTEGER PRIMARY KEY AUTOINCREMENT, 
   postID INTEGER REFERENCES people(userID), 
	userID INTEGER REFERENCES people(userID),
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


func (data *Forum) GetUser() []User {
	user := []User{}

	rows, _ := data.DB.Query(`
	SELECT * FROM people
	`)

	var uID int 
	var uuid string
	var username string
	var email string 
	var pass string 


	for rows.Next(){
		rows.Scan(&uID,&uuid,&username,&email,&pass)
		user =append(user, User{
			UserID: uID,
			Uuid: uuid,
			Username: username,
			Email: email,
			Password: pass,
		})
	}
	return user
}



func (data *Forum) GetPost() []PostFeed {

	posts := []PostFeed{}

	rows, _ := data.DB.Query(`
	SELECT * FROM post
`)

	var id int
	var username string
	var title string
	var content string
	var likes int
	var dislikes int
	var created string
	var category string

	for rows.Next() {
		rows.Scan(&id, &username, &title, &content, &likes, &dislikes, &category, &created)

		posts = append(posts, PostFeed{
			PostID:    id,
			Username:  username,
			Title:     title,
			Content:   content,
			Likes:     likes,
			Dislikes:  dislikes,
			Category:  category,
			CreatedAt: created,
		})
	}
	//fmt.Println(posts)
	return posts
}

func (data *Forum) GetComments() []Comment {
	comments := []Comment{}
	rows, _ := data.DB.Query(`SELECT * FROM comments`)

	var commentid int
	var postid int
	var userid int
	var content string
	var created string

	for rows.Next() {
		rows.Scan(&commentid, &postid, &userid, &content, &created)
		comments = append(comments, Comment{
			CommentID: commentid,
			PostID:    postid,
			UserId:    userid,
			Content:   content,
			CreatedAt: created,
		})
	}
	return comments
}

func (data *Forum) GetSession() []Session {

	session := []Session{}

	rows, _ := data.DB.Query(`
 SELECT * FROM session
 `)

   var session_token string 
	var uName string
	var exTime time.Time

	for rows.Next() {
		rows.Scan(&session_token,&uName, &exTime)
		session = append(session, Session{
			SessionID: session_token,
			Username: uName,
			Expiry:   exTime,
		})
	}
	return session

}

