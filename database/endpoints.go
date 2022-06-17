package database

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (p PostFeed) MarshallJSON() ([]byte, error) {
	return json.Marshal(p)
}

// creates all needed templates
// will need to be reduced as there is too many at the moment.
// var tpl *template.Template

// parses files for all templates allowing them to be called.
// func init() {
// 	tpl = template.Must(template.ParseGlob("templates/*"))
// }

// sessions
// var sessions = map[string]Session{}

// func (s Session) isExpired() bool {
// 	return s.Expiry.Before(time.Now())

// }

// @TODO: error handling.
// login page.
func (data *Forum) LoginWeb(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginUser is running********")
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	// switch r.Method {
	// case "POST":
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("LoginWeb ParseForm error: %+v\n", err)
		return
	}

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
	err = row.Scan(&passwordHash)

	tpl := template.Must(template.ParseGlob("templates/*"))
	if err != nil {
		err := tpl.ExecuteTemplate(w, "login.html", "check username and password")
		if err != nil {
			fmt.Printf("LoginWeb ExecuteTemplate error: %+v\n", err)
			return
		}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	// returns nil on succcess
	if err == nil {
		err = data.CreateSession(Session{
			SessionID: sessionToken.String(),
			Username:  user.Username,
			Expiry:    expiresAt,
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken.String(),
			Expires: expiresAt,
			// MaxAge:  2 * int(time.Hour),
		})
		fmt.Println(data.GetSession())
		// w.WriteHeader(200)
		http.Redirect(w, r, "/home", http.StatusFound)
		// data.HomePage(w, r)
	} else {
		fmt.Println("incorrect password")
		err := tpl.ExecuteTemplate(w, "login.html", "check username and password")
		if err != nil {
			fmt.Printf("LoginWeb ExecuteTemplate error: %+v\n", err)
			return
		}
	}
}

// @TODO: error handling.
func (data *Forum) GetSignupPage(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
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
	tpl := template.Must(template.ParseGlob("templates/*"))
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
		err := tpl.ExecuteTemplate(w, "signup.html", "please check username and/or password criteria")
		if err != nil {
			fmt.Printf("SignUpUser ExecuteTemplate signup.html error: %+v\n", err)
			return
		}
		return
	}

	row := data.DB.QueryRow("SELECT uuid FROM people where username =?", user.Username)
	var username string
	err := row.Scan(&username)
	if err != sql.ErrNoRows {
		fmt.Printf("sql scan row user error: %+v\n", err)
		err1 := tpl.ExecuteTemplate(w, "signup.html", "username taken")
		if err1 != nil {
			fmt.Printf("SignUpUser ExecuteTemplate (username) error1: %+v\n", err1)
			return
		}
	}
	row = data.DB.QueryRow("SELECT uuid FROM people where email =?", user.Email)
	var userEmail string
	err = row.Scan(&userEmail)
	if err != sql.ErrNoRows {
		fmt.Printf("sql scan row email error: %+v\n", err)
		err2 := tpl.ExecuteTemplate(w, "signup.html", "e-mail in use")
		if err2 != nil {
			fmt.Printf("SignUpUser ExecuteTemplate (email) error2: %+v\n", err2)
			return
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		err3 := tpl.ExecuteTemplate(w, "signup.html", "there was an error registering account")
		if err3 != nil {
			fmt.Printf("SignUpUser ExecuteTemplate (password) error:  %+v\n", err3)
			return
		}
		return
	}

	sessionID := uuid.NewV4()
	err = data.CreateUser(User{
		Uuid:     sessionID.String(),
		Username: user.Username,
		Email:    user.Email,
		Password: string(passwordHash),
	})

	if err != nil {
		err4 := tpl.ExecuteTemplate(w, "signup.html", "there was an error registering account")
		if err4 != nil {
			fmt.Printf("SignUpUser ExecuteTemplate (password) error:  %+v\n", err4)
			return
		}
	}
	http.Redirect(w, r, "/login", 302)
	return
}

// check cookie
func (data *Forum) CheckCookie(writer http.ResponseWriter, request *http.Request) bool {
	c, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Printf("CheckCookie (Cookie) error: %+v\n", err)
			return false
		}
	}

	sessionToken := c.Value
	// var currentSession Session
	a, err := data.GetSession()
	if err != nil {
		fmt.Printf("CheckCookie (GetSession) error: %+v\n", err)
	}

	// fmt.Println(sessionToken == a[0].SessionID)
	// sessFound := false

	for _, sess := range a {
		fmt.Println(sessionToken, " : ", sess.SessionID)
		if sessionToken == sess.SessionID {
			// fmt.Println(sessionToken, " : ", sess.SessionID)
			// currentSession = sess
			// sessFound = true
		}

		// if !sessFound {
		// // // 	//writer.WriteHeader(http.StatusUnauthorized)
		// // 	return
		// // }

		// if currentSession.isExpired() {
		// 	data.DB.Exec("DELETE FROM session where sessionID ='" + currentSession.SessionID + "'")
		// }
	}
	return true
}

func (data *Forum) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		fmt.Printf("Logout Cookie error: %+v\n", err)
	}

	sessionToken := c.Value
	var currentSession Session
	a, err := data.GetSession()
	if err != nil {
		fmt.Printf("Logout GetSession error: %+v\n", err)
	}

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

// home page.
func (data *Forum) HomePage(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	tpl := template.Must(template.ParseGlob("templates/*"))

	if err := request.ParseForm(); err != nil { // checks for errors parsing form
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("HomePage (ParseForm) error:  %+v\n", err)
		return
	}
	loggedIn := data.CheckCookie(writer, request)
	// 🐈
	if !loggedIn {
		data, err := data.GetPost()
		if err != nil {
			fmt.Printf("HomePage (GetPost) error: %+v\n", err)
		}
		err = tpl.ExecuteTemplate(writer, "guest.html", data)
		if err != nil {
			fmt.Printf("HomePage ExecuteTemplate (guest.html) error: %+v\n", err)
		}
		return

	} else {
		post, _ := data.GetPost()
		lastPost := post[len(post)-1]

		postCategory := request.FormValue("category")
		postTitle := request.FormValue("title")
		postContent := request.FormValue("content")

		postLikes := 0
		postDislikes := 0
		time := time.Now()
		postCreated := time.Format("01-02-2006 15:04")

		// checks session and selects the last one (the latest one)
		sess, _ := data.GetSession()
		currentSession := sess[len(sess)-1]

		user := currentSession.Username // fetches username from session

		type postSessionStruct struct {
			Post        []PostFeed
			UserSession Session
		}

		var postAndSession postSessionStruct

		postAndSession.UserSession = currentSession

		// checks if last post == current submit values to prevent duplicate posts
		if lastPost.Content == postContent {
			fmt.Println("duplicate")
			postAndSession.Post, _ = data.GetPost()
			tpl.ExecuteTemplate(writer, "./home", postAndSession)

		} else {
			// postAndSession.UserSession = data.GetSession()[0]

			if postTitle != "" || postContent != "" || postCategory != "" {

				err := data.CreatePost(PostFeed{
					// User:      sessionID.String(),

					Username:  user,
					Title:     postTitle,
					Content:   postContent,
					Likes:     postLikes,
					Dislikes:  postDislikes,
					Category:  postCategory,
					CreatedAt: postCreated,
				})
				if err != nil {
					fmt.Printf("HomePage (CreatePost) items error: %+v\n", err)
					return
				}

				postAndSession.Post, err = data.GetPost()
				if err != nil {
					fmt.Printf("HomePage (GetPost) items error: %+v\n", err)
					return
				}

				err = tpl.ExecuteTemplate(writer, "./home", postAndSession)
				if err != nil {
					fmt.Printf("HomePage ExecuteTemplate user homepage error: %+v\n", err)
					return
				}

			}
		}
		data, err := data.GetPost()
		postAndSession.Post = data
		if err != nil {
			fmt.Printf("HomePage (GetPost) data error: %+v\n", err)
			return
		}
		err = tpl.ExecuteTemplate(writer, "home.html", postAndSession)
		if err != nil {
			fmt.Printf("HomePage ExecuteTemplate (home.html) error: %+v\n", err)
			return
		}
	}
}

func (data *Forum) GuestView(writer http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	items, err := data.GetPost()
	if err != nil {
		fmt.Printf("GuestView (GetPost) items error: %+v\n", err)
		return
	}
	err = tpl.ExecuteTemplate(writer, "guest.html", items)
	if err != nil {
		fmt.Printf("GuestView ExecuteTemplate error: %+v\n", err)
		return
	}
}

func (data *Forum) CategoriesList(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	loggedIn := data.CheckCookie(w, r)
	if !loggedIn {
		err := tpl.ExecuteTemplate(w, "guestCategories.html", nil)
		if err != nil {
			fmt.Printf("CategoriesList ExecuteTemplate error: %+v\n", err)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(w, "categories.html", nil)
	if err != nil {
		fmt.Printf("CategoriesList ExecuteTemplate error: %+v\n", err)
		return
	}
}

func (data *Forum) CategoryDump(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("CategoryDump (ParseForm) error: %+v\n", err)
		return
	}

	loggedIn := data.CheckCookie(w, r)

	type CategoryPost struct { // create a []post in order to store multiple posts
		Post []PostFeed
	}

	var postByCategory CategoryPost // create variable to link to our struct
	category := r.URL.Path
	cat := ""
	if !loggedIn {
		cat = strings.Replace(category, "/categoryg/", "", -1) // we use replace instead of trim as we are working with strings-- and useful characters were being removed
	} else {
		cat = strings.Replace(category, "/category/", "", -1) // we use replace instead of trim as we are working with strings-- and useful characters were being removed
	}

	posts, err := data.GetPost()
	if err != nil {
		fmt.Printf("CategoryDump (GetPost) posts error: %+v\n", err)
		return
	} // get all posts
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
		err := tpl.ExecuteTemplate(w, "guestCategoryPosts.html", postByCategory)
		if err != nil {
			fmt.Printf("CategoryDump ExecuteTemplate (guestCategoryPosts.html) error: %+v\n", err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	err1 := tpl.ExecuteTemplate(w, "categoryPosts.html", postByCategory)
	if err1 != nil {
		fmt.Printf("CategoryDump ExecuteTemplate (categoryPosts.html) error: %+v\n", err1)
	}
}

func (data *Forum) PwReset(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(writer, "passwordReset.html", nil)
}

func (data *Forum) UserProfile(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")

	type profile struct {
		Profile UsrProfile
		// UserSession Session
	}
	sess, _ := data.GetSession()
	currentSession := sess[len(sess)-1]
	// data.GetSession()[len(data.GetSession())-1]

	var User profile

	// User.UserSession =currentSession

	User.Profile.Name = currentSession.Username

	User.Profile.Info = "Hello my name is panda and I like to sleep and eat bamboo--- nom"
	User.Profile.Gender = "Panda"
	User.Profile.Age = 7
	User.Profile.Location = "Bamboo Forest"

	err := tpl.ExecuteTemplate(writer, "profile.html", User)
	if err != nil {
		fmt.Printf("UserProfile ExecuteTemplate (profile.html) error: %+v\n", err)
		return
	}
}

// Threads handles posts and their comments-- and displays them on /threads.
func (data *Forum) Threads(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	w.WriteHeader(http.StatusOK)
	// grab current url, parse the form to allow taking data from html
	url := r.URL.Path
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Threads ParseForm error: %+v\n", err)
		return
	}
	idstr := strings.Trim(url, "/thread") // trim text so  we are only left with the final end point (postID)
	// fmt.Println(idstr)
	id, err := strconv.Atoi(idstr) // convert to number as postID is stored as an int on our database
	if err != nil {
		http.Error(w, "400 Bad Request", 400)
	}

	comment := r.FormValue("comment") // take "comment" id value from html form
	time := time.Now()                // create a new time variable using following format
	postCreated := time.Format("01-02-2006 15:04")
	var postWithComments Databases

	post, err := data.GetPost() // get all posts
	if err != nil {
		fmt.Printf("Threads (GetPost) posts error: %+v\n", err)
		return
	}
	sess, _ := data.GetSession()
	currentSession := sess[len(sess)-1]
	// data.GetSession()[len(data.GetSession())-1]
	cmnt, _ := data.GetComments()
	lastComment := cmnt[len(cmnt)-1]
	// data.GetComments()[len(data.GetComments())-1]

	// if last comment != current submitted values then create a comment, otherwise show comments
	if lastComment.Content != comment {
		// if comment from html is not an empty string, add a new value to our comment database using the following structure
		if comment != "" || comment == " " {
			err = data.CreateComment(Comment{
				PostID:    post[id-1].PostID, // id-1 is used as items on database start at index 0, but start at 1 on html url
				UserId:    currentSession.Username,
				Content:   comment,
				CreatedAt: postCreated,
			})
			if err != nil {
				fmt.Printf("Threads (CreateComment) error: %+v\n", err)
				return
			}
		}
	}

	if id > len(post) { // checks so that a post that is not higher than total post amount and returns an error
		http.Error(w, "404 post not found", 404)
	}
	commentdb, err := data.GetComments() // get data from comment database
	if err != nil {
		fmt.Printf("Threads (GetComments) error: %+v\n", err)
		return
	}
	// only adds a comment into database if the post id matches the url id (post requested)--- to only fetch the same ids
	for _, comment := range commentdb {
		if comment.PostID == id {
			postWithComments.Comment = append(postWithComments.Comment, comment) // only adds matching comments to the database to be called only for specific posts
		}
	}
	postWithComments.Post = post[id-1] // only allows us to send the requested post

	err = tpl.ExecuteTemplate(w, "thread.html", postWithComments)
	if err != nil {
		fmt.Printf("Threads ExecuteTemplate (thread.html) error: %+v\n", err)
		return
	}
}

func (data *Forum) ThreadGuest(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	w.WriteHeader(http.StatusOK)
	// grab current url, parse the form to allow taking data from html
	url := r.URL.Path
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Threads ParseForm error: %+v\n", err)
		return
	}

	idstr := strings.Trim(url, "/threadg") // trim text so  we are only left with the final end point (postID)
	id, err := strconv.Atoi(idstr)         // convert to number as postID is stored as an int on our database
	if err != nil {
		http.Error(w, "400 Bad Request", 400)
	}

	var postWithComments Databases
	post, err := data.GetPost() // get all posts
	if err != nil {
		fmt.Printf("ThreadGuest (GetPost) posts error: %+v\n\n", err)
		return
	}
	if id > len(post) { // checks so that a post that is not higher than total post amount and returns an error
		http.Error(w, "404 post not found", 400)
	}
	commentdb, err := data.GetComments() // get data from comment database
	if err != nil {
		fmt.Printf("ThreadGuest (GetComments) error: %+v\n", err)
		return
	}
	// only adds a comment into database if the post id matches the url id (post requested)--- to only fetch the same ids
	for _, comment := range commentdb {
		// fmt.Println("value", v, "comment ", comment)
		if comment.PostID == id {
			postWithComments.Comment = append(postWithComments.Comment, comment) // only adds matching comments to the database to be called only for specific posts
			// fmt.Println(comment)
		}
	}

	postWithComments.Post = post[id-1] // only allows us to send the requested post
	err = tpl.ExecuteTemplate(w, "threadGuest.html", postWithComments)
	if err != nil {
		fmt.Printf("ThreadGuest ExecuteTemplate (threadGuest.html) error: %+v\n", err)
		return
	}
}

func (data *Forum) AboutFunc(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	loggedIn := data.CheckCookie(w, r)
	if !loggedIn {
		err := tpl.ExecuteTemplate(w, "aboutGuest.html", nil)
		if err != nil {
			fmt.Printf("AboutFunc ExecuteTemplate (aboutGuest.html) error: %+v\n\n", err)
			return
		}
	} else {
		err := tpl.ExecuteTemplate(w, "about.html", nil)
		if err != nil {
			fmt.Printf("AboutFunc ExecuteTemplate (about.html) error: %+v\n", err)
			return
		}
	}
}

func (data *Forum) ContactUs(w http.ResponseWriter, r *http.Request) {
	loggedIn := data.CheckCookie(w, r)
	tpl := template.Must(template.ParseGlob("templates/*"))
	if !loggedIn {
		err := tpl.ExecuteTemplate(w, "contactGuest.html", nil)
		if err != nil {
			fmt.Printf("ThreadGuest ExecuteTemplate (threadGuest.html) error: %+v\n", err)
			return
		}
	} else {
		err := tpl.ExecuteTemplate(w, "contact-us.html", nil)
		if err != nil {
			fmt.Printf("ContactUs ExecuteTemplate (contact-us.html) error: %+v\n", err)
			return
		}
	}
}

func (data *Forum) UserPhoto(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "photo.html", nil)
	if err != nil {
		fmt.Printf("UserPhoto ExecuteTemplate (photo.html) error: %+v\n", err)
		return
	}
}

func (data *Forum) UserPosts(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")

	user, _ := data.GetSession()
	currentUser := user[len(user)-1]
	// if user.session == user in post --- send this post

	posts, _ := data.GetPost()

	type UserPosts struct {
		Post []PostFeed
	}
	var usrPosts UserPosts

	for _, post := range posts {
		if post.Username == currentUser.Username {
			usrPosts.Post = append(usrPosts.Post, post)
			// fmt.Println(currentUser.Username, post.Username)
		}
	}
	err := tpl.ExecuteTemplate(writer, "posts.html", usrPosts)
	if err != nil {
		fmt.Printf("UserPosts ExecuteTemplate (posts.html) error: %+v\n", err)
		return
	}
}

func (data *Forum) UserComments(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "comments.html", nil)
	if err != nil {
		fmt.Printf("UserComments ExecuteTemplate (comments.html) error: %+v\n", err)
		return
	}
}

func (data *Forum) UserLikes(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "likes.html", nil)
	if err != nil {
		fmt.Printf("UserLikes ExecuteTemplate (likes.html) error: %+v\n", err)
	}
}

func (data *Forum) UserDislikes(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "likes.html", nil)
	if err != nil {
		fmt.Printf("UserDislikes ExecuteTemplate (likes.html) error: %+v\n", err)
	}
}

func (data *Forum) UserShares(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "shares.html", nil)
	if err != nil {
		fmt.Printf("UserShares ExecuteTemplate (shares.html) error: %+v\n", err)
		return
	}
}

func (data *Forum) UserInfo(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "userinfo.html", nil)
	if err != nil {
		fmt.Printf("UserInfo ExecuteTemplate (userinfo.html) error: %+v\n", err)
		return
	}
}

func (data *Forum) Customization(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	err := tpl.ExecuteTemplate(writer, "customize.html", nil)
	if err != nil {
		fmt.Printf("Customization ExecuteTemplate (customize.html) error: %+v\n", err)
		return
	}
}

func (data *Forum) HandleLikeDislike(writer http.ResponseWriter, request *http.Request, isLike bool) {
	items, err := data.GetPost()
	if err != nil {
		fmt.Printf("HandleLikeDislike (GetPost) posts error: %+v\n", err)
		return
	}
	reqItemIDraw := request.URL.Query().Get("id")
	reqItemID, err := strconv.Atoi(reqItemIDraw)
	if err != nil {
		fmt.Printf("unable to parse post id: %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte("{\"500\": \"Error parsing post id\"}"))
		if err != nil {
			fmt.Printf("unable to send json response for post %d\n", reqItemID)
			return
		}
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
		_, err := writer.Write([]byte("{\"404\": \"Error finding post\"}"))
		if err != nil {
			fmt.Printf("unable to send json response for post %d\n", reqItemID)
			return
		}
		return
	}

	j, err := requestedItem.MarshallJSON()
	if err != nil {
		fmt.Printf("unable to marshal json for post %d: %v\n", reqItemID, err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte("{\"500\": \"Error marshalling json for post\"}"))
		if err != nil {
			fmt.Printf("unable to send json response for post %d\n", reqItemID)
			return
		}
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
		if isLike {
			requestedItem.Likes = requestedItem.Likes + 1
		} else {
			requestedItem.Dislikes = requestedItem.Dislikes + 1
		}

		err := data.UpdatePost(requestedItem)
		if err != nil {
			fmt.Printf("unable to modify dis-likes for post %d: %v\n", reqItemID, err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Header().Set("Content-Type", "application/json")
			_, err := writer.Write([]byte("{\"500\": \"Error modifying dis-likes for post\"}"))
			if err != nil {
				fmt.Printf("unable to send json response for post %d\n", reqItemID)
				return
			}
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write([]byte("{\"success\":true}"))
		if err != nil {
			fmt.Printf("unable to send json response for post %d\n", reqItemID)
			return
		}
		fmt.Printf("modified dis-likes on post %d\n", reqItemID)
	}
}

func (data *Forum) ImgUpload(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.html")
		t.Execute(writer, token)
	} else {
		err := request.ParseMultipartForm(5 << 20)
		if err != nil {
			http.Error(writer, "The uploaded image is too big. Please use an image less than 20mb in size", http.StatusBadRequest)
			return
		}
		file, handler, err := request.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(writer, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
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
		data.GuestView(w, r)

		// user handlers
	case "/photo":
		data.UserPhoto(w, r)
	case "/posts":
		data.UserPosts(w, r)
	case "/upload":
		data.ImgUpload(w, r)
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
		data.HandleLikeDislike(w, r, true)
	case "/dislike":
		data.HandleLikeDislike(w, r, false)
	}
}
