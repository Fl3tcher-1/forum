package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

var DB *sql.DB

// @TODO: handle errors for all create funcs.
func (forum *Forum) CreateUser(user User) error {
	stmt, err := forum.DB.Prepare("INSERT INTO people (uuid, username, email, password) VALUES (?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("CreateUser DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(user.Uuid, user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("CreateUser Exec error: %+v\n", err)
	}
	defer stmt.Close()
	return nil
}

func (forum *Forum) CreateSession(session Session) error {
	stmt, err := forum.DB.Prepare("INSERT INTO session (sessionID, userName, expiryTime) VALUES (?, ?, ?);")
	if err != nil {
		return fmt.Errorf("CreateSession DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(session.SessionID, session.Username, session.Expiry)
	if err != nil {
		return fmt.Errorf("CreateSession Exec error: %+v\n", err)
	}
	defer stmt.Close()
	return nil
}

func (forum *Forum) CreatePost(post PostFeed) error {
	stmt, err := forum.DB.Prepare("INSERT INTO post (username, title, content, likes, dislikes, category, dateCreated, image) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("CreatePost DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(post.Username, post.Title, post.Content, post.Likes, post.Dislikes, post.Category, post.CreatedAt, post.Image)
	if err != nil {
		return fmt.Errorf("CreatePost Exec error: %+v\n", err)
	}
	defer stmt.Close()
	return nil
}

func (forum *Forum) CreateComment(comment Comment) error {
	stmt, err := forum.DB.Prepare("INSERT INTO comments ( postID, userID, content, dateCreated) VALUES (?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("CreateComment DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(comment.PostID, comment.UserId, comment.Content, comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("CreateComment Exec error: %+v\n", err)
	}
	defer stmt.Close()
	return nil
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

func userTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS people (
	userID INTEGER PRIMARY KEY,	
	uuid TEXT, 
	username TEXT UNIQUE,
	email TEXT UNIQUE, 
	password TEXT);
`)
	if err != nil {
		return fmt.Errorf("userTable DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("userTable Exec error: %+v\n", err)
	}
	return nil
}

func sessionTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS session (
	sessionID TEXT PRIMARY KEY REFERENCES people(uuid),	
	userName TEXT REFERENCES people(username), 
	expiryTime TEXT);
	`)
	if err != nil {
		return fmt.Errorf("sessionTable DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("sessionTable Exec error: %+v\n", err)
	}
	return nil
}

func postTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "post" (
 "postID" INTEGER PRIMARY KEY AUTOINCREMENT,
 "username" TEXT REFERENCES session(userName),
 "title" TEXT,
 "content" TEXT, 
 "likes" INTEGER,
 "dislikes" INTEGER,
 "category" TEXT,
 "dateCreated" TEXT,
 "image" TEXT
	);
 `)
	if err != nil {
		return fmt.Errorf("postTable DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("postTable Exec error: %+v\n", err)
	}
	return nil
}

// @TODO: add likes/dislikes to comments.
func commentTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments (
   commentID INTEGER PRIMARY KEY AUTOINCREMENT, 
   postID INTEGER REFERENCES people(userID), 
	userID STRING REFERENCES session(userName),
	content TEXT NOT NULL, 
	dateCreated TEXT NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("commentTable DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("commentTable Exec error: %+v\n", err)
	}
	return nil
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

func (data *Forum) GetPost() ([]PostFeed, error) {
	posts := []PostFeed{}
	rows, err := data.DB.Query(`SELECT * FROM post`)
	if err != nil {
		return posts, fmt.Errorf("GetPost DB Query error: %+v\n", err)
	}
	var id int
	var uiD string
	var title string
	var content string
	var likes int
	var dislikes int
	var created string
	var category string
	var image interface{}

	for rows.Next() {
		err := rows.Scan(&id, &uiD, &title, &content, &likes, &dislikes, &category, &created)
		if err != nil {
			return posts, fmt.Errorf("GetPost rows.Scan error: %+v\n", err)
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
			Image:     image,
		})
	}
	// fmt.Println(posts)
	return posts, nil
}

// @TODO: add likes/dislikes to comments.
func (data *Forum) GetComments() ([]Comment, error) {
	comments := []Comment{}
	rows, err := data.DB.Query(`SELECT * FROM comments`)
	if err != nil {
		return comments, fmt.Errorf("GetComments DB Query error: %+v\n", err)
	}
	var commentid int
	var postid int
	var userid string
	var content string
	var created string

	for rows.Next() {
		err := rows.Scan(&commentid, &postid, &userid, &content, &created)
		if err != nil {
			return comments, fmt.Errorf("GetComments rows.Scan error: %+v\n", err)
		}
		comments = append(comments, Comment{
			CommentID: commentid,
			PostID:    postid,
			UserId:    userid,
			Content:   content,
			CreatedAt: created,
		})
	}
	return comments, nil
}

func (data *Forum) GetSession() ([]Session, error) {
	session := []Session{}
	rows, err := data.DB.Query(`SELECT * FROM session`)
	if err != nil {
		return session, fmt.Errorf("GetSession DB Query error: %+v\n", err)
	}
	var session_token string
	var uName string
	var exTime string

	for rows.Next() {
		err := rows.Scan(&session_token, &uName, &exTime)
		if err != nil {
			return nil, fmt.Errorf("GetSession rows.Scan error: %+v\n", err)
		}
		// time.Format("01-02-2006 15:04")
		convTime, err := time.Parse("2006-01-02 15:04:05.999999999Z07:00", exTime)
		if err != nil {
			return nil, fmt.Errorf("GetSession time.Parse(exTime) error: %+v\n", err)
		}
		session = append(session, Session{
			SessionID: session_token,
			Username:  uName,
			Expiry:    convTime,
		})
	}
	return session, nil
}
