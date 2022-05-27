package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	// uuid "v2/go/pkg/mod/github.com/satori/go.uuid@v1.2.0"
	// "v2/go/pkg/mod/golang.org/x/crypto@v0.0.0-20200622213623-75b288015ac9/bcrypt"
)

type Log struct {
	Loggedin bool
}

type User struct {
	Username string
	Password string
	Email    string
	UserID   string
}

// could it be used to store data for userprofile and use a single template execution???
// holds details of user session-- used for cookies
type session struct {
	Id     int
	Uuid   string // random value to be stored at the browser
	Email  string
	UserID string
	// CreatedAt	time.Time
}

type usrProfile struct {
	Name string
	// image    *os.Open
	Info     string
	Photo    string
	Gender   string
	Age      int
	Location string
	Posts    []string
	Comments []string
	Likes    []Reaction
	Dislikes []Reaction
	Shares   []string
	Userinfo map[string]string
	// custom   string
}

type Post struct {
	Title    string
	Content  string
	Date     string
	Comments int
	PostID   string
	UserID   string
	Reaction Reaction
}

type Reaction struct {
	PostID     string
	UserID     string
	ReactionID string
	CommentID  string
	// React      int
	Likes    int
	Dislikes int
}

type PostFeed struct {
	ID       int    `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Content  string `json:"content,omitempty"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
	Created  string `json:"created,omitempty"`
	Category string `json:"category,omitempty"`
	// Comments []string
}

func (p PostFeed) MarshallJSON() ([]byte, error) {
	return json.Marshal(p)
}

// creates all needed templates
// will need to be reduced as there is too many at the moment
var tpl *template.Template

// parses files for all templates allowing them to be called
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

// login page
func LoginWeb(w http.ResponseWriter, r *http.Request) {
	var Roles []string
	var registered Log
	var passwordHash string
	var user User

	Roles = append(Roles, "guest", "user", "moderator", "admin")
	registered.Loggedin = false

	fmt.Println(registered.Loggedin)
	fmt.Println(Roles)

	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	fmt.Println(r.FormValue("username"))
	fmt.Println(r.FormValue("username"))
	
	r.ParseForm()

	stmt := "SELECT passwordHash FROM people WHERE Username = ?"
	row := DB.QueryRow(stmt, user.Username)
	err := row.Scan(&passwordHash)
	if err != nil {
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	// returns nil on succcess
	if err == nil {
		tpl.ExecuteTemplate(w, "home.html", nil)
		http.Redirect(w, r, "/home.html", http.StatusFound)
		return
	}

	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "login.html", "check username and password")

	w.WriteHeader(http.StatusOK)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "500 Internal Server Error", 500)
		fmt.Printf("LoginWeb(writeheader) error:  %+v\n", err)
	}
	tpl.ExecuteTemplate(w, "login.html", nil)
}

func GetSignupPage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "signup.html", nil)
}

/*  1. check e-mail criteria
    2. check u.username criteria
	 3. check password criteria
	 4. check if u.username is already exists in database
	 5. create bcrypt hash from password
	 6. insert u.username and password hash in database
*/
func SignUpUser(w http.ResponseWriter, r *http.Request) {
	var user User

	r.ParseForm() // parses sign up form to fetch needed information
	user.Email = r.FormValue("email")
	// check if e-mail is valid format
	isValidEmail := true

	if isValidEmail != strings.Contains(user.Email, "@") || isValidEmail != strings.Contains(user.Email, ".") { // checks if e-mail is valid by checking if it contains "@"
		isValidEmail = false
	}

	user.Username = r.FormValue("username")
	// check if only alphanumerical numbers
	isAlphaNumeric := true

	for _, char := range user.Username {
		if unicode.IsLetter(char) && unicode.IsNumber(char) { // checks if character not a special character
			isAlphaNumeric = false
		}
	}
	// checks if name length meets criteria
	nameLength := (5 <= len(user.Username) && len(user.Username) <= 50)

	fmt.Println(nameLength)

	// check pw criteria
	user.Password = r.FormValue("password")

	fmt.Println(user)
	var pwLower, pwUpper, pwNumber, pwSpace, pwLength bool
	pwSpace = false

	for _, char := range user.Password {
		switch {
		case unicode.IsLower(char):
			pwLower = true
		case unicode.IsUpper(char):
			pwUpper = true
		case unicode.IsNumber(char):
			pwNumber = true
		// case unicode.IsPunct(char) || unicode.IsSymbol(char):
		// 	pwSpecial = true
		case unicode.IsSpace(int32(char)):
			pwSpace = true
		}
	}

	minPwLength := 8
	maxPwLength := 30

	if minPwLength <= len(user.Password) && len(user.Password) <= maxPwLength {
		pwLength = true
	}

	if !pwLower || !pwUpper || !pwNumber || !pwLength || pwSpace || !isAlphaNumeric || !nameLength {
		tpl.ExecuteTemplate(w, "signup.html", "please check usrname and/or password criteria")
		return
	}

	stmt := "SELECT id FROM people where username = ?"

	row := DB.QueryRow(stmt, user.Username)

	var id string
	err := row.Scan(&id)
	if err != sql.ErrNoRows {
		// fmt.Println("user exists", err)
		tpl.ExecuteTemplate(w, "sign-up.html", "username taken")
		fmt.Printf("sql scan row id error: %+v\n", err)
		return
	}
	stmt = "SELECT id FROM people where email =?"

	row = DB.QueryRow(stmt, user.Username)

	var userEmail string
	err = row.Scan(&userEmail)
	if err != sql.ErrNoRows {
		fmt.Printf("sql scan row email error: %+v\n", err)
		tpl.ExecuteTemplate(w, "signup.html", "e-mail in use")
	}

	var passwordHash []byte

	passwordHash, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		tpl.ExecuteTemplate(w, "signup.html", "there was an error registering account")
		fmt.Printf("Register Account (passwordHash) error:  %+v\n", err)
		return
	}

	var insertStmt *sql.Stmt
	insertStmt, err = DB.Prepare("INSERT INTO people (username, email, passwordHASH) VALUES (?, ?, ?);")
	if err != nil {
		tpl.ExecuteTemplate(w, "signup.html", "there was an error registering account")
		fmt.Printf("Register Account (insertStmt) error:  %+v\n", err)

		return
	}
	defer insertStmt.Close()

	var result sql.Result
	result, err = insertStmt.Exec(user.Username, user.Email, passwordHash)
	rowsAff, err1 := result.RowsAffected()
	if err1 != nil {
		fmt.Printf("rowsAff: %+v error:  %+v\n", rowsAff, err1)
	}
	lastIns, err2 := result.LastInsertId()
	if err2 != nil {
		fmt.Printf("lastIns: %+v error:  %+v\n", lastIns, err2)
	}
	if err != nil {
		tpl.ExecuteTemplate(w, "signup.html", "there was an error registering account")
		fmt.Printf("Register Account (result) error:  %+v\n", err)
		return
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// home page
func HomePage(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	if err := request.ParseForm(); err != nil { // checks for errors parsing form
		http.Error(writer, "500 Internal Server Error", 500)
		fmt.Printf("ParseForm (HomePage) error:  %+v\n", err)
		return
	}
	// 🐈
	UserDatabase()

	posts, err := sql.Open("sqlite3", "./database/feed.db")
	if err != nil {
		fmt.Printf("posts sql.Open (HomePage) error:  %+v\n", err)
	}
	feed := Feed(posts)

	items := feed.Get()
	postStuff := request.ParseForm()
	fmt.Println(postStuff)

	postCategory := request.FormValue("category")
	postTitle := request.FormValue("title")
	postContent := request.FormValue("content")
	postLikes := 0
	// fetch from database and add one to likes and dislikes
	postDislikes := 0
	time := time.Now()
	postCreated := time.Format("01-02-2006 15:04")

	// check to see if title, content and category has been provided to stop making empty posts
	if postTitle != "" || postContent != "" || postCategory != "" {

		// add values into database
		feed.Add(PostFeed{
			Title:    postTitle,
			Content:  postContent,
			Likes:    postLikes,
			Dislikes: postDislikes,
			Created:  postCreated,
			Category: postCategory,
		})

		tpl.ExecuteTemplate(writer, "./home", items)
	}

	tpl.ExecuteTemplate(writer, "home.html", items)
}

func CategoriesList(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(writer, "categories.html", nil)
}

func PwReset(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(writer, "passwordReset.html", nil)
}

func UserProfile(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")

	var users User

	users.Username = "test"

	var usrInfo usrProfile

	usrInfo.Name = "Panda"
	usrInfo.Info = "Hello my name is panda and I like to sleep and eat bamboo--- nom"
	usrInfo.Gender = "Panda"
	usrInfo.Age = 7
	usrInfo.Location = "Bamboo Forest"

	tpl.ExecuteTemplate(writer, "profile.html", usrInfo)
}

func Threads(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var postInfo Post
	postInfo.Title = "testing"
	postInfo.Content = "this is a completely empty post"
	postInfo.Comments = 2
	postInfo.Date = "11/11/11"

	fmt.Print(postInfo)
	tpl.ExecuteTemplate(w, "thread.html", postInfo)
}

func AboutFunc(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "about.html", nil)
	if err != nil {
		fmt.Printf("AboutFunc Execute.Template error: %+v\n", err)
	}
}

func ContactUs(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "contact-us.html", nil)
	if err != nil {
		fmt.Printf("ContactUs Execute.Template error: %+v\n", err)
	}
}

func UserPhoto(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "photo.html", nil)
	if err != nil {
		fmt.Printf("UserPhoto Execute.Template error: %+v\n", err)
	}
}

func UserPosts(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "posts.html", nil)
	if err != nil {
		fmt.Printf("UserPosts Execute.Template error: %+v\n", err)
	}
}

func UserComments(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "comments.html", nil)
	if err != nil {
		fmt.Printf("UserComments Execute.Template error: %+v\n", err)
	}
}

func UserLikes(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "likesanddislikes.html", nil)
	if err != nil {
		fmt.Printf("UserLikes Execute.Template error: %+v\n", err)
	}
}

func UserDislikes(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "likesanddislikes.html", nil)
	if err != nil {
		fmt.Printf("UserDislikes Execute.Template error: %+v\n", err)
	}
}

func UserShares(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "shares.html", nil)
	if err != nil {
		fmt.Printf("UserShares Execute.Template error: %+v\n", err)
	}
}

func UserInfo(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "userinfo.html", nil)
	if err != nil {
		fmt.Printf("UserInfo Execute.Template error: %+v\n", err)
	}
}

func Customization(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "customize.html", nil)
	if err != nil {
		fmt.Printf("Customization Execute.Template error: %+v\n", err)
	}
}

func AddLike(writer http.ResponseWriter, request *http.Request) {
	UserDatabase()

	posts, err := sql.Open("sqlite3", "./database/feed.db")
	if err != nil {
		fmt.Printf("posts sql.Open error:  %+v\n", err)
	}
	feed := Feed(posts)

	items := feed.Get()
	reqItemIDraw := request.URL.Query().Get("id")
	reqItemID, err := strconv.Atoi(reqItemIDraw)
	if err != nil {
		fmt.Printf("unable to parse post id: %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte("{\"500\": \"Error parsing post id\"}"))
		return
	}
	requestedItem := PostFeed{}

	for _, item := range items {
		if item.ID == reqItemID {
			requestedItem = item
		}
	}

	if requestedItem.Created == "" {
		fmt.Printf("unable to find post %d in db: %v\n", reqItemID, err)
		writer.WriteHeader(http.StatusNotFound)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte("{\"404\": \"Error finding post\"}"))
		return
	}

	j, err := requestedItem.MarshallJSON()
	if err != nil {
		fmt.Printf("unable to marshal json for post %d: %v\n", reqItemID, err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte("{\"500\": \"Error marshalling json for post\"}"))
		return
	}

	switch request.Method {
	case http.MethodGet:
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write(j)
		if err != nil {
			fmt.Printf("unable to send json response for post %d\n", reqItemID)
		}
	case http.MethodPost:

		requestedItem.Likes = requestedItem.Likes + 1

		err := feed.Update(requestedItem)
		if err != nil {
			fmt.Printf("unable increment likes for post %d: %v\n", reqItemID, err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Header().Set("Content-Type", "application/json")
			writer.Write([]byte("{\"500\": \"Error incrementing likes for post\"}"))
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write([]byte("{\"success\":true}"))
		if err != nil {
			fmt.Printf("unable to send json response for post %d\n", reqItemID)
		}
		fmt.Printf("added like to post %d\n", reqItemID)
	}
}

func AddDislike(writer http.ResponseWriter, request *http.Request) {
	UserDatabase()

	posts, err := sql.Open("sqlite3", "./database/feed.db")
	if err != nil {
		fmt.Printf("posts sql.Open error:  %+v\n", err)
	}
	feed := Feed(posts)

	items := feed.Get()
	reqItemIDraw := request.URL.Query().Get("id")
	reqItemID, err := strconv.Atoi(reqItemIDraw)
	if err != nil {
		fmt.Printf("unable to parse post id: %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte("{\"500\": \"Error parsing post id\"}"))
		return
	}
	requestedItem := PostFeed{}

	for _, item := range items {
		if item.ID == reqItemID {
			requestedItem = item
		}
	}

	if requestedItem.Created == "" {
		fmt.Printf("unable to find post %d in db: %v\n", reqItemID, err)
		writer.WriteHeader(http.StatusNotFound)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte("{\"404\": \"Error finding post\"}"))
		return
	}

	j, err := requestedItem.MarshallJSON()
	if err != nil {
		fmt.Printf("unable to marshal json for post %d: %v\n", reqItemID, err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte("{\"500\": \"Error marshalling json for post\"}"))
		return
	}

	switch request.Method {
	case http.MethodGet:
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write(j)
		if err != nil {
			fmt.Printf("unable to send json response for post %d\n", reqItemID)
		}
	case http.MethodPost:
		requestedItem.Dislikes = requestedItem.Dislikes + 1
		err := feed.Update(requestedItem)
		if err != nil {
			fmt.Printf("unable increment dislikes for post %d: %v\n", reqItemID, err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Header().Set("Content-Type", "application/json")
			writer.Write([]byte("{\"500\": \"Error incrementing dislikes for post\"}"))
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write([]byte("{\"success\":true}"))
		if err != nil {
			fmt.Printf("unable to send json response for post %d\n", reqItemID)
		}
		fmt.Printf("added dislike to post %d\n", reqItemID)
	}
}
