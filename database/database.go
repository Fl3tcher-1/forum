package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

var DB *sql.DB

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
	stmt, err := forum.DB.Prepare("INSERT INTO post (username, title, content, category, dateCreated) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("CreatePost DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(post.Username, post.Title, post.Content, post.Category, post.CreatedAt)
	if err != nil {
		return fmt.Errorf("CreatePost Exec error: %+v\n", err)
	}
	defer stmt.Close()
	return nil
}

func (forum *Forum) CreateReaction(reaction Reaction) error {
	stmt, err := forum.DB.Prepare("INSERT INTO reaction (postid, userid, reactionid, commentid, liked, disliked) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("CreateReaction DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(reaction.PostID, reaction.UserID, reaction.ReactionID, reaction.CommentID, reaction.Liked, reaction.Disliked)
	if err != nil {
		return fmt.Errorf("CreateReactions Exec error: %+v\n", err)
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
	stmt, err := feed.DB.Prepare("UPDATE post SET title = ?, content = ?, category = ? WHERE postID = ?;")
	if err != nil {
		return fmt.Errorf("UpdatePost DB Prepare error: %+v", err)
	}
	defer stmt.Close()
	// stmt.QueryRow(stmt, item.Title, item.Content, item.Category)
	_, err = stmt.Exec(item.Title, item.Content, item.Category, item.PostID)
	if err != nil {
		return fmt.Errorf("unable to insert item into post: %w", err)
	}
	return nil
}

func (feed *Forum) UpdateReaction(item Reaction) error {
	stmt, err := feed.DB.Prepare("UPDATE reaction SET liked = ?, disliked = ?;")
	if err != nil {
		return fmt.Errorf("UpdateReaction DB Prepare error: %+v", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(item.Liked, item.Disliked)
	if err != nil {
		return fmt.Errorf("unable to update reaction: %w", err)
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
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS post (
 postID INTEGER PRIMARY KEY AUTOINCREMENT,
 username TEXT REFERENCES session(userName),
 title TEXT,
 content TEXT, 
 category TEXT,
 dateCreated TEXT);
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

func reactionsTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS reaction (
   reactionID INTEGER PRIMARY KEY AUTOINCREMENT,
   postID INTEGER REFERENCES posts(postID),
   userID INTEGER REFERENCES people(userID),
   commentID INTEGER REFERENCES comments(commentID),
   liked BOOL,
   disliked BOOL);
	`)
	if err != nil {
		return fmt.Errorf("reactionsTable DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("reactionsTable Exec error: %+v\n", err)
	}
	return nil
}

func Connect(db *sql.DB) *Forum {
	userTable(db)
	sessionTable(db)
	postTable(db)
	commentTable(db)
	reactionsTable(db)

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
	var created string
	var category string

	for rows.Next() {
		err := rows.Scan(&id, &uiD, &title, &content, &category, &created)
		if err != nil {
			return posts, fmt.Errorf("GetPost rows.Scan error: %+v\n", err)
		}

		likes, err := getLikesForPost(data.DB, id)
		if err != nil {
			return posts, fmt.Errorf("GetPost getLikesForPost error: %+v\n", err)
		}

		dislikes, err := getDislikesForPost(data.DB, id)
		if err != nil {
			return posts, fmt.Errorf("GetPost getDislikesForPost error: %+v\n", err)
		}

		posts = append(posts, PostFeed{
			PostID:    id,
			Username:  uiD,
			Title:     title,
			Content:   content,
			Category:  category,
			CreatedAt: created,
			Likes:     likes,
			Dislikes:  dislikes,
		})
	}

	// fmt.Println(posts)
	return posts, nil
}

func getLikesForPost(db *sql.DB, id int) (int, error) {
	stmt, err := db.Prepare("SELECT liked FROM reaction WHERE liked = TRUE AND postID = ?")
	if err != nil {
		return 0, fmt.Errorf("getLikesForPost DB Prepare error: %+v\n", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return 0, fmt.Errorf("getLikesForPost DB Query error: %+v\n", err)
	}

	counter := 0
	for rows.Next() {
		counter++
	}
	return counter, nil
}

func getDislikesForPost(db *sql.DB, id int) (int, error) {
	stmt, err := db.Prepare("SELECT disliked FROM reaction WHERE disliked = TRUE AND postID = ?")
	if err != nil {
		return 0, fmt.Errorf("getDislikesForPost DB Prepare error: %+v\n", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return 0, fmt.Errorf("getDislikesForPost DB Query error: %+v\n", err)
	}

	counter := 0
	for rows.Next() {
		counter++
	}
	return counter, nil
}

// TODO: implement the get reaction(post, comment, etc)
func (data *Forum) GetReactions() ([]Reaction, error) {
	reactions := []Reaction{}
	rows, err := data.DB.Query(`SELECT * FROM reaction`)
	if err != nil {
		return reactions, fmt.Errorf("GetReactions DB Query error: %+v\n", err)
	}

	var reactionID int
	var postID int
	var userID int
	var commentID int
	var liked bool
	var disliked bool

	for rows.Next() {
		err := rows.Scan(&reactionID, &postID, &userID, &commentID, &liked, &disliked)
		if err != nil {
			return reactions, fmt.Errorf("GetReactions rows.Scan error: %+v\n", err)
		}
		reactions = append(reactions, Reaction{
			ReactionID: reactionID,
			PostID:     postID,
			UserID:     userID,
			CommentID:  commentID,
			Liked:      liked,
			Disliked:   disliked,
		})
	}
	fmt.Println(reactions)
	return reactions, nil
}

func (data *Forum) GetReactionByPostID(targetPostID, targetUserID string) (Reaction, error) {
	stmt, err := data.DB.Prepare("SELECT * FROM reaction WHERE postID = ? AND userID = ?")
	if err != nil {
		return Reaction{}, fmt.Errorf("GetReactionByPostID DB Prepare error: %+v", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(targetPostID, targetUserID)
	if err != nil {
		return Reaction{}, fmt.Errorf("GetReactionByPostID DB Query error: %+v", err)
	}

	var reactionID int
	var postID int
	var userID int
	var commentID int
	var liked bool
	var disliked bool

	for rows.Next() {
		err := rows.Scan(&reactionID, &postID, &userID, &commentID, &liked, &disliked)
		if err != nil {
			return Reaction{}, fmt.Errorf("GetReactionByPostID rows.Scan error: %+v\n", err)
		}
		return Reaction{
			ReactionID: reactionID,
			PostID:     postID,
			UserID:     userID,
			CommentID:  commentID,
			Liked:      liked,
			Disliked:   disliked,
		}, nil
	}
	return Reaction{}, nil
}

// @TODO: add likes/dislikes(reactions) to comments.
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
