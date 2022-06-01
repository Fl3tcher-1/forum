package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"html/template"
	"net/http"
	"strings"
	"unicode"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Log struct {
	Loggedin bool
}

// type User struct {
// 	Username string
// 	Password string
// 	Email    string
// }

// could it be used to store data for userprofile and use a single template execution???

// holds details of user session-- used for cookies

type Post struct {
	Title    string
	Content  string
	Date     string
	Comments int
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
	Likes    []string
	Shares   []string
	Userinfo map[string]string
	// custom   string
}

// creates all needed templates
// will need to be reduced as there is too many at the moment

var tpl *template.Template

// parses files for all templates allowing them to be called
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

// sessions

// var sessions = map[string]Session{}

func (s Session) isExpired() bool {
	return s.Expiry.Before(time.Now())

}

// login page
func (data *Forum) LoginWeb(w http.ResponseWriter, r *http.Request) {

	fmt.Println("*****loginUser is running********")

	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	r.ParseForm()

	var user User

	sessionToken := uuid.NewV4()
	expiresAt := time.Now().Add(120 * time.Second)

	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")

	data.CreateSession(Session{
		SessionID: sessionToken.String(),
		Username:  user.Username,
		Expiry:    expiresAt,
	})

	var passwordHash string

	row := data.DB.QueryRow("SELECT password FROM people WHERE Username = ?", user.Username)
	err := row.Scan(&passwordHash)

	if err != nil {
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	// returns nill on succcess
	if err == nil {
		// posts, err := sql.Open("sqlite3", "./database/feed.db")
		// if err != nil {
		// 	database.CheckErr(err)
		// }
		// feed := database.Feed(posts)
		// items := feed.Get()
		// registered.Loggedin = true
		// fmt.Println(registered)
		// tpl.ExecuteTemplate(w, "home.html", items)
		http.Redirect(w, r, "/home", 302)
		return
	}

	// fmt.Println("incorrect password")
	// tpl.ExecuteTemplate(w, "login.html", "check username and password")
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	// returns nill on succcess
	if err == nil {

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken.String(),
			Expires: expiresAt,
		})

		// var cookie *http.Cookie
		// cookie, err = r.Cookie("session")
		// if err != nil {
		// 	sID := uuid.NewV4()
		// 	//fmt.Println("cookie was not found")
		// 	cookie = &http.Cookie{
		// 	Name:  "session",
		// 	Value: sID.String(),
		// 	//Secure:   true,
		// 	HttpOnly: true,
		// 	MaxAge: 2 * int(time.Hour),
		// 	}
		// 	http.SetCookie(w, cookie)
		// 	//w.WriteHeader(200)
		//tpl.ExecuteTemplate(w, "home.html", nil)

		http.Redirect(w, r, "/home", 302)
	} else {
		fmt.Println("incorrect password")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	fmt.Println("here")

	//tpl.ExecuteTemplate(w, "login.html", nil)
}

func (data *Forum) GetSignupPage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "signup.html", nil)
}

/*  1. check e-mail criteria
    2. check u.username criteria
	 3. check password criteria
	 4. check if u.username is already exists in database
	 5. create bcrypt hash from password
	 6. insert u.username and password hash in database
*/
func (data *Forum) SignUpUser(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() // parses sign up form to fetch needed information

	fmt.Println("****Sign-up new user is running ")

	var user User

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

	row := data.DB.QueryRow("SELECT uuid FROM people where username = ?", user.Username)
	var username string
	err := row.Scan(&username)
	if err != sql.ErrNoRows {
		// fmt.Println("user exists", err)
		tpl.ExecuteTemplate(w, "signup.html", "username taken")
		fmt.Printf("sql scan row id error: %+v\n", err)
		return
	}
	row = data.DB.QueryRow("SELECT uuid FROM people where email =?", user.Email)
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

	sessionID := uuid.NewV4()

	data.CreateUser(User{
		Uuid:     sessionID.String(),
		Username: user.Username,
		Email:    user.Email,
		Password: string(passwordHash),
	})

	if err != nil {
		tpl.ExecuteTemplate(w, "signup.html", "there was an error registering account")
		//fmt.Printf("Register Account (insertStmt) error:  %+v\n", err)
		//defer data.Close()
		return
	} else {
		http.Redirect(w, r, "/login", 302)
		return
	}
}

// home page
func (data *Forum) HomePage(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	if err := request.ParseForm(); err != nil { // checks for errors parsing form
		http.Error(writer, "500 Internal Server Error", 500)
		fmt.Printf("ParseForm (HomePage) error:  %+v\n", err)
		return
	}

	// 🐈

	postCategory := request.FormValue("category")

	postTitle := request.FormValue("title")

	postContent := request.FormValue("content")
	postLikes := 0
	postDislikes := 1
	time := time.Now()
	postCreated := time.Format("01-02-2006 15:04")

	sessionID := uuid.NewV4()
	user := 1

	fmt.Println(postCategory)
	fmt.Println(postTitle)
	fmt.Println(postContent)

	if postTitle != "" || postContent != "" || postCategory != "" {

		data.CreatePost(PostFeed{
			Uuid:      sessionID.String(),
			Title:     postTitle,
			Content:   postContent,
			Likes:     postLikes,
			Dislikes:  postDislikes,
			Category:  postCategory,
			CreatedAt: postCreated,
			UserID:    user,
		},
			User{
				UserID: user,
			})

		items := data.Get()
		fmt.Println(items)

		tpl.ExecuteTemplate(writer, "./home", items)
	}

	tpl.ExecuteTemplate(writer, "home.html", data.Get())

}

func (data *Forum) CategoriesList(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(writer, "categories.html", nil)
}

func (data *Forum) PwReset(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(writer, "passwordReset.html", nil)
}

func (data *Forum) UserProfile(writer http.ResponseWriter, request *http.Request) {
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

//Threds handles posts and their comments-- and displays them on /threads
func (data *Forum) Threads(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	//grab current url, parse the form to allow taking data from html 
	url:= r.URL.Path
	r.ParseForm()

	idstr := strings.Trim(url, "/thread") //trim text so  we are only left with the final end point (postID)
	// fmt.Println(idstr)

	id, err := strconv.Atoi(idstr) //convert to number as postID is stored as an int on our database
	if err != nil {
		http.Error(w, "400 Bad Request", 400)
	}

	comment := r.FormValue("comment") //take "comment" id value from html form
	time := time.Now() //create a new time variable using following format
	postCreated := time.Format("01-02-2006 15:04")

	//Databases holds our post and comment databases
	type Databases struct {
		Post    PostFeed
		Comment []Comment
	}

	var postWithComments Databases

	post := data.Get() // get all posts

	//if comment from html is not an empty string, add a new value to our comment database using the following structure
	if comment != "" { 
		data.CreateComment(Comment{
			PostID:    post[id-1].PostID, //id-1 is used as items on database start at index 0, but start at 1 on html url
			UserId:    post[0].PostID,
			Content:   comment,
			CreatedAt: postCreated,
		})
	}
	if id > len(post){ //checks so that a post that is not higher than total post amount and returns an error
		http.Error(w, "404 post not found", 400)
		return
	}
	commentdb := data.GetComments() // get data from comment database

	//only adds a comment into database if the post id matches the url id (post requested)--- to only fetch the same ids
	for _, comment := range commentdb{
		// fmt.Println("value", v, "comment ", comment)
		if comment.PostID == id{
			postWithComments.Comment = append(postWithComments.Comment, comment) //only adds matching comments to the database to be called only for specific posts
			// fmt.Println(comment)
		}
	}

	postWithComments.Post = post[id-1] //only allows us to send the requested post

	tpl.ExecuteTemplate(w, "thread.html", postWithComments)

}

func (data *Forum) AboutFunc(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "about.html", nil)
	if err != nil {
		fmt.Printf("AboutFunc Execute.Template error: %+v\n", err)
	}
}

func (data *Forum) ContactUs(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "contact-us.html", nil)
	if err != nil {
		fmt.Printf("ContactUs Execute.Template error: %+v\n", err)
	}
}

func (data *Forum) UserPhoto(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "photo.html", nil)
	if err != nil {
		fmt.Printf("UserPhoto Execute.Template error: %+v\n", err)
	}
}

func (data *Forum) UserPosts(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "posts.html", nil)
	if err != nil {
		fmt.Printf("UserPosts Execute.Template error: %+v\n", err)
	}
}

func (data *Forum) UserComments(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "comments.html", nil)
	if err != nil {
		fmt.Printf("UserComments Execute.Template error: %+v\n", err)
	}
}

func (data *Forum) UserLikes(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "likes.html", nil)
	if err != nil {
		fmt.Printf("UserLikes Execute.Template error: %+v\n", err)
	}
}

func (data *Forum) UserShares(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "shares.html", nil)
	if err != nil {
		fmt.Printf("UserShares Execute.Template error: %+v\n", err)
	}
}

func (data *Forum) UserInfo(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "userinfo.html", nil)
	if err != nil {
		fmt.Printf("UserInfo Execute.Template error: %+v\n", err)
	}
}

func (data *Forum) Customization(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "customize.html", nil)
	if err != nil {
		fmt.Printf("Customization Execute.Template error: %+v\n", err)
	}
}

func (data *Forum) Handler(w http.ResponseWriter, r *http.Request) {
	// check for cookie
	switch r.URL.Path {
	// page handlers
	case "/stylesheet": // handle css
		http.ServeFile(w, r, "./templates/stylesheet.css")
	case "/":
		data.LoginWeb(w, r)
	case "/login":
		data.LoginWeb(w, r)
	case "/home":
		data.HomePage(w, r)
	case "/categories":
		data.CategoriesList(w, r)
	case "/reset":
		data.PwReset(w, r)
	case "/signup":
		data.GetSignupPage(w, r)
	case "/sign-up-form":
		data.SignUpUser(w, r)
	case "/profile":
		data.UserProfile(w, r)
	// case "/thread/*":
	// 	data.Threads(w, r)
	case "/about":
		data.AboutFunc(w, r)
	case "/contact-us":
		data.ContactUs(w, r)

		// user handlers
	case "/photo":
		data.UserPhoto(w, r)
	case "/posts":
		data.UserPosts(w, r)
	case "/comments":
		data.UserComments(w, r)
	case "/likes":
		data.UserLikes(w, r)
	case "/shares":
		data.UserShares(w, r)
	case "/info":
		data.UserInfo(w, r)
	case "/custom":
		data.Customization(w, r)

		// handles images
	case "/cat":
		http.ServeFile(w, r, "./images/cat.jpg")
	case "/chicken":
		http.ServeFile(w, r, "./images/chicken.jpeg")
	case "/cow":
		http.ServeFile(w, r, "./images/cow.jpg")
	case "/hamster":
		http.ServeFile(w, r, "./images/hamster.jpg")
	case "/owl":
		http.ServeFile(w, r, "./images/owl.jpg")
	case "/panda":
		http.ServeFile(w, r, "./images/panda.jpg")
	case "/shark":
		http.ServeFile(w, r, "./images/shark.jpg")
	case "/doge":
		http.ServeFile(w, r, "./images/doge.jpg")
	case "/question":
		http.ServeFile(w, r, "./images/question.jpg")
	}
}

// func logOut(){

// close session
// log user out
// clear cookie
// }
