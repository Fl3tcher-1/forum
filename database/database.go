package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

type Forum struct {
	*sql.DB
}

var DB *sql.DB

// func CheckErr(err error) {
// 	fmt.Println(err)
// 	log.Fatal(err)
// }

// @TODO: handle errors for all create funcs.
func (forum *Forum) CreateUser(user User) {
	stmt, err := forum.DB.Prepare("INSERT INTO people (uuid, username, email, password) VALUES (?, ?, ?, ?);")
	if err != nil {
		fmt.Printf("CreateUser DB Prepare error: %+v\n", err)
	}
	stmt.Exec(user.Uuid, user.Username, user.Email, user.Password)
	defer stmt.Close()
}

func (forum *Forum) CreateSession(session Session) {
	stmt, err := forum.DB.Prepare("INSERT INTO session (sessionID, userName, expiryTime) VALUES (?, ?, ?);")
	if err != nil {
		fmt.Printf("CreateSession DB Prepare error: %+v\n", err)
	}
	stmt.Exec(session.SessionID, session.Username, session.Expiry)
	defer stmt.Close()
}

func (forum *Forum) CreatePost(post PostFeed) {
	stmt, err := forum.DB.Prepare("INSERT INTO post (username, title, content, likes, dislikes, category, dateCreated) VALUES (?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		fmt.Printf("CreatePost DB Prepare error: %+v\n", err)
	}
	stmt.Exec(post.Username, post.Title, post.Content, post.Likes, post.Dislikes, post.Category, post.CreatedAt)
	defer stmt.Close()
}

func (forum *Forum) CreateComment(comment Comment) {
	stmt, err := forum.DB.Prepare("INSERT INTO comments ( postID, userID, content, dateCreated) VALUES (?, ?, ?, ?);")
	if err != nil {
		fmt.Printf("CreateComment DB Prepare error: %+v\n", err)
	}
	stmt.Exec(comment.PostID, comment.UserId, comment.Content, comment.CreatedAt)
	defer stmt.Close()

}

// Update(Updates an item in a table).
func (feed *Forum) UpdatePost(item PostFeed) error {
	stmt, err := feed.DB.Prepare("UPDATE post SET title = ?, content = ?, likes = ?, dislikes = ?, category = ? WHERE postID = ?;")
	if err != nil {
		return fmt.Errorf("UpdatePost DB Prepare error: %+v", err)
	}
	defer stmt.Close()
	// stmt.QueryRow(stmt, item.Title, item.Content, item.Category)
	_, err = stmt.Exec(item.Title, item.Content, item.Likes, item.Dislikes, item.Category, item.PostID)
	if err != nil {
		return fmt.Errorf("unable to insert item into post: %w", err)
	}
	return nil
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
		fmt.Printf("userTable DB Prepare error: %+v\n", err)
	}
	stmt.Exec()
}

func sessionTable(db *sql.DB) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS session (
	sessionID TEXT PRIMARY KEY REFERENCES people(uuid),	
	userName TEXT REFERENCES people(username), 
	expiryTime TEXT);
	`)
	if err != nil {
		fmt.Printf("sessionTable DB Prepare error: %+v\n", err)
	}
	stmt.Exec()
}

func postTable(db *sql.DB) {
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
		fmt.Printf("postTable DB Prepare error: %+v\n", err)
	}
	stmt.Exec()
}

// @TODO: add likes/dislikes to comments.
func commentTable(db *sql.DB) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments (
   commentID INTEGER PRIMARY KEY AUTOINCREMENT, 
   postID INTEGER REFERENCES people(userID), 
	userID INTEGER REFERENCES people(userID),
	content TEXT NOT NULL, 
	dateCreated TEXT NOT NULL);
	`)
	if err != nil {
		fmt.Printf("commentTable DB Prepare error: %+v\n", err)
	}
	stmt.Exec()
}

func Connect(db *sql.DB) *Forum {
	userTable(db)
	sessionTable(db)
	postTable(db)
	commentTable(db)

	return &Forum{
		DB: db,
	}
}

func (data *Forum) GetPost() []PostFeed {
	posts := []PostFeed{}
	rows, err := data.DB.Query(`SELECT * FROM post`)
	if err != nil {
		fmt.Printf("GetPost DB Query error: %+v\n", err)
	}
	var id int
	var uiD string
	var title string
	var content string
	var likes int
	var dislikes int
	var created string
	var category string

	for rows.Next() {
		err := rows.Scan(&id, &uiD, &title, &content, &likes, &dislikes, &category, &created)
		if err != nil {
			fmt.Printf("GetPost rows.Scan error: %+v\n", err)
		}
		posts = append(posts, PostFeed{
			PostID:    id,
			Username:  uiD,
			Title:     title,
			Content:   content,
			Likes:     likes,
			Dislikes:  dislikes,
			Category:  category,
			CreatedAt: created,
		})
	}
	// fmt.Println(posts)
	return posts
}

// @TODO: add likes/dislikes to comments.
func (data *Forum) GetComments() []Comment {
	comments := []Comment{}
	rows, err := data.DB.Query(`SELECT * FROM comments`)
	if err != nil {
		fmt.Printf("GetComments DB Query error: %+v\n", err)
	}
	var commentid int
	var postid int
	var userid int
	var content string
	var created string

	for rows.Next() {
		err := rows.Scan(&commentid, &postid, &userid, &content, &created)
		if err != nil {
			fmt.Printf("GetComments rows.Scan error: %+v\n", err)
		}
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
	rows, err := data.DB.Query(`SELECT * FROM session`)
	if err != nil {
		fmt.Printf("GetSession DB Query error: %+v\n", err)
	}
	var session_token string
	var uName string
	var exTime time.Time

	for rows.Next() {
		err := rows.Scan(&session_token, &uName, &exTime)
		if err != nil {
			fmt.Printf("GetSession rows.Scan error: %+v\n", err)
		}
		session = append(session, Session{
			SessionID: session_token,
			Username:  uName,
			Expiry:    exTime,
		})
	}

	return session
}
