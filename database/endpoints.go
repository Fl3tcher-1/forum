package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"html/template"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// type Log struct {
// 	Loggedin bool

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

// sessions

// var sessions = map[string]Session{}

func (s Session) isExpired() bool {
	return s.Expiry.Before(time.Now())

}

// @TODO: error handling
// login page
func (data *Forum) LoginWeb(w http.ResponseWriter, r *http.Request) {

	fmt.Println("*****loginUser is running********")

	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// switch r.Method {
	// case "POST":
	r.ParseForm()

	var user User
	// sessionToken := uuid.NewV4()
	// expiresAt := time.Now().Add(120 * time.Second)

	// user.Username = r.FormValue("username")
	// user.Password = r.FormValue("password")

	// data.CreateSession(Session{
	// 	SessionID: sessionToken.String(),
	// 	Username:  user.Username,
	// 	Expiry:    expiresAt,
	// })

	sessionToken := uuid.NewV4()
	expiresAt := time.Now().Add(120 * time.Second)

	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")

	var passwordHash string

	row := data.DB.QueryRow("SELECT password FROM people WHERE Username = ?", user.Username)
	err := row.Scan(&passwordHash)

	if err != nil {
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	// returns nil on succcess
	if err == nil {

		data.CreateSession(Session{
			SessionID: sessionToken.String(),
			Username:  user.Username,
			Expiry:    expiresAt,
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken.String(),
			Expires: expiresAt,
			//MaxAge:  2 * int(time.Hour),
		})
		//w.WriteHeader(200)
		http.Redirect(w, r, "/home", 302)
		//data.HomePage(w, r)
	} else {
		fmt.Println("incorrect password")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
	}

}

// @TODO: error handling
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
		return
	} else {
		http.Redirect(w, r, "/login", 302)
		return
	}
}

// check cookie

func (data *Forum) CheckCookie(writer http.ResponseWriter, request *http.Request) bool {

	c, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println(err)
			return false
		}
	}

	sessionToken := c.Value
	var currentSession Session
	a := data.GetSession()

	fmt.Println(sessionToken)
	//fmt.Println(sessionToken == a[0].SessionID)
	//sessFound := false

	for _, sess := range a {
		fmt.Println(sessionToken, " : ", sess.SessionID)
		if sessionToken == sess.SessionID {
			//fmt.Println(sessionToken, " : ", sess.SessionID)
			currentSession = sess
			//sessFound = true
		}

		// if !sessFound {
		// // // 	//writer.WriteHeader(http.StatusUnauthorized)
		// // 	return
		// // }

		if currentSession.isExpired() {
			data.DB.Exec("DELETE FROM session where sessionID ='" + currentSession.SessionID + "'")
		}
	}
	return true

}

func (data *Forum) Logout(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session_token")

	sessionToken := c.Value
	var currentSession Session
	a := data.GetSession()

	fmt.Println(sessionToken)

	fmt.Println("here")

	for _, sess := range a {
		if sessionToken == sess.SessionID {
			currentSession = sess
			data.DB.Exec("DELETE FROM session where sessionID ='" + currentSession.SessionID + "'")
		}
	}

	c = &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

// home page
func (data *Forum) HomePage(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	if err := request.ParseForm(); err != nil { // checks for errors parsing form
		http.Error(writer, "500 Internal Server Error", 500)
		fmt.Printf("ParseForm (HomePage) error:  %+v\n", err)
		return
	}

	loggedIn := data.CheckCookie(writer, request)
	fmt.Println(loggedIn)

	// 🐈
	if !loggedIn {
		fmt.Println(loggedIn)
		err := tpl.ExecuteTemplate(writer, "guest.html", data.GetPost()) 
		if err != nil {
			fmt.Printf("ExecuteTemplate guest error: %+v", err)
		}
		return

	} else {

		postCategory := request.FormValue("category")

		postTitle := request.FormValue("title")

		postContent := request.FormValue("content")
		postLikes := 0
		postDislikes := 0
		time := time.Now()
		postCreated := time.Format("01-02-2006 15:04")

		user := "1"

		fmt.Println(postCategory)
		fmt.Println(postTitle)
		fmt.Println(postContent)

		if postTitle != "" || postContent != "" || postCategory != "" {

			data.CreatePost(PostFeed{
				//User:      sessionID.String(),

				Username:  user,
				Title:     postTitle,
				Content:   postContent,
				Likes:     postLikes,
				Dislikes:  postDislikes,
				Category:  postCategory,
				CreatedAt: postCreated,
			})

			items := data.GetPost()
			fmt.Println(items)

			tpl.ExecuteTemplate(writer, "./home", items)
		}
		tpl.ExecuteTemplate(writer, "home.html", data.GetPost())
	}
}

func (data *Forum) Guestview(writer http.ResponseWriter, r *http.Request) {

	fmt.Println("here")
	items := data.GetPost()
	fmt.Println(data.GetPost())
	tpl.ExecuteTemplate(writer, "guest.html", items)
}

func (data *Forum) CategoriesList(w http.ResponseWriter, r *http.Request) {

	loggedIn := data.CheckCookie(w, r)

	if !loggedIn {
		tpl.ExecuteTemplate(w, "guestCategories.html", nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "categories.html", nil)
}

func (data *Forum) CategoryDump(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	loggedIn := data.CheckCookie(w, r)

	type CategoryPost struct { // create a []post in order to store multiple posts
		Post []PostFeed
	}
	var postByCategory CategoryPost //create variable to link to our struct

	category := r.URL.Path
	cat := ""
	if !loggedIn {
		cat = strings.Replace(category, "/categoryg/", "", -1) //we use replace instead of trim as we are working with strings-- and useful characters were being removed
	} else {
		cat = strings.Replace(category, "/category/", "", -1) //we use replace instead of trim as we are working with strings-- and useful characters were being removed
	}

	posts := data.GetPost() // get all posts
	// fmt.Println(posts)
	// check every post to find ones whose category matches our url path
	categoryFound := false // used to check if a valid category was entered
	for _, post := range posts {
		// fmt.Println(cat, post.Category)
		fmt.Println(post.Category)
		if cat == post.Category {
			// fmt.Println(post)
			categoryFound = true
			postByCategory.Post = append(postByCategory.Post, post) // add the matching post to our post[] in struct
		}
	}
	if !categoryFound {
		http.Error(w, "404 category not found or has no posts", 404)
		return
	}

	if !loggedIn {
		tpl.ExecuteTemplate(w, "guestCategoryPosts.html", postByCategory)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "categoryPosts.html", postByCategory)

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

	var usrInfo UsrProfile

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
	url := r.URL.Path
	r.ParseForm()

	idstr := strings.Trim(url, "/thread") //trim text so  we are only left with the final end point (postID)
	// fmt.Println(idstr)

	id, err := strconv.Atoi(idstr) //convert to number as postID is stored as an int on our database
	if err != nil {
		http.Error(w, "400 Bad Request", 400)
	}

	comment := r.FormValue("comment") //take "comment" id value from html form
	time := time.Now()                //create a new time variable using following format
	postCreated := time.Format("01-02-2006 15:04")

	//Databases holds our post and comment databases
	type Databases struct {
		Post    PostFeed
		Comment []Comment
	}

	var postWithComments Databases

	post := data.GetPost() // get all posts

	//if comment from html is not an empty string, add a new value to our comment database using the following structure
	if comment != "" || comment == " " {
		data.CreateComment(Comment{
			PostID:    post[id-1].PostID, //id-1 is used as items on database start at index 0, but start at 1 on html url
			UserId:    post[0].PostID,
			Content:   comment,
			CreatedAt: postCreated,
		})
	}
	if id > len(post) { //checks so that a post that is not higher than total post amount and returns an error
		http.Error(w, "404 post not found", 404)
		return
	}
	commentdb := data.GetComments() // get data from comment database

	//only adds a comment into database if the post id matches the url id (post requested)--- to only fetch the same ids
	for _, comment := range commentdb {
		// fmt.Println("value", v, "comment ", comment)
		if comment.PostID == id {
			postWithComments.Comment = append(postWithComments.Comment, comment) //only adds matching comments to the database to be called only for specific posts
			// fmt.Println(comment)
		}
	}

	postWithComments.Post = post[id-1] //only allows us to send the requested post

	tpl.ExecuteTemplate(w, "thread.html", postWithComments)

}
func (data *Forum) ThreadGuest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	//grab current url, parse the form to allow taking data from html
	url := r.URL.Path
	r.ParseForm()

	idstr := strings.Trim(url, "/threadg") //trim text so  we are only left with the final end point (postID)
	// fmt.Println(idstr)

	id, err := strconv.Atoi(idstr) //convert to number as postID is stored as an int on our database
	if err != nil {
		http.Error(w, "400 Bad Request", 400)
	}

	type Databases struct {
		Post    PostFeed
		Comment []Comment
	}

	var postWithComments Databases

	post := data.GetPost() // get all posts

	if id > len(post) { //checks so that a post that is not higher than total post amount and returns an error
		http.Error(w, "404 post not found", 400)
		return
	}
	commentdb := data.GetComments() // get data from comment database

	//only adds a comment into database if the post id matches the url id (post requested)--- to only fetch the same ids
	for _, comment := range commentdb {
		// fmt.Println("value", v, "comment ", comment)
		if comment.PostID == id {
			postWithComments.Comment = append(postWithComments.Comment, comment) //only adds matching comments to the database to be called only for specific posts
			// fmt.Println(comment)
		}
	}

	postWithComments.Post = post[id-1] //only allows us to send the requested post

	tpl.ExecuteTemplate(w, "threadGuest.html", postWithComments)

	// fmt.Println("here")
	// items:= data.GetPost()
	// fmt.Println(data.GetPost())
	// tpl.ExecuteTemplate(w,"threadGuest.html", items)
}

func (data *Forum) AboutFunc(w http.ResponseWriter, r *http.Request) {
	loggedIn := data.CheckCookie(w, r)
	if !loggedIn {
		tpl.ExecuteTemplate(w, "aboutGuest.html", nil)
		return

	} else {
		err := tpl.ExecuteTemplate(w, "about.html", nil)
		if err != nil {
			fmt.Printf("AboutFunc Execute.Template error: %+v\n", err)
		}
	}
}

func (data *Forum) ContactUs(w http.ResponseWriter, r *http.Request) {
	loggedIn := data.CheckCookie(w, r)

	if !loggedIn {
		tpl.ExecuteTemplate(w, "contactGuest.html", nil)
		return

	} else {
		err := tpl.ExecuteTemplate(w, "contact-us.html", nil)
		if err != nil {
			fmt.Printf("ContactUs Execute.Template error: %+v\n", err)
		}
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

func (data *Forum) UserDislikes(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "likes.html", nil)
	if err != nil {
		fmt.Printf("UserDislikes Execute.Template error: %+v\n", err)
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

func (data *Forum) AddLike(writer http.ResponseWriter, request *http.Request) {
	items := data.GetPost()
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
		if item.PostID == reqItemID {
			requestedItem = item
		}
	}

	if requestedItem.CreatedAt == "" {
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

		err := data.UpdatePost(requestedItem)
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

func (data *Forum) AddDislike(writer http.ResponseWriter, request *http.Request) {
	items := data.GetPost()
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
		if item.PostID == reqItemID {
			requestedItem = item
		}
	}

	if requestedItem.CreatedAt == "" {
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
		err := data.UpdatePost(requestedItem)
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

func (data *Forum) Handler(w http.ResponseWriter, r *http.Request) {

	// data.CheckCookie(w, r)

	switch r.URL.Path {
	// page handlers
	case "/stylesheet": // handle css
		http.ServeFile(w, r, "./templates/stylesheet.css")
	case "/":
		data.LoginWeb(w, r)
	case "/login":
		data.LoginWeb(w, r)
	case "/logout":
		data.Logout(w, r)
	case "/home":
		data.HomePage(w, r)
	case "/categories":
		data.CategoriesList(w, r)
	case "/guestcategories":
		data.CategoriesList(w, r)
	case "/reset":
		data.PwReset(w, r)
	case "/signup":
		data.GetSignupPage(w, r)
	case "/sign-up-form":
		data.SignUpUser(w, r)
	case "/profile":
		data.UserProfile(w, r)
	case "/thread":
		data.Threads(w, r)
		// data.UserProfile(w, r)
	// case "/thread/*":
	// 	data.Threads(w, r)
	case "/about":
		data.AboutFunc(w, r)
	case "/contact-us":
		data.ContactUs(w, r)
	case "/guest":
		data.Guestview(w, r)

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

		// api handlers
	case "/like":
		data.AddLike(w, r)
	case "/dislike":
		data.AddDislike(w, r)
	}
}
