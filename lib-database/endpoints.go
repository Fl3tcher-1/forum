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

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (p PostFeed) MarshallJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (c Comment) MarshallJSON() ([]byte, error) {
	return json.Marshal(c)
}

// @TODO: error handling.
// login page.
func (data *Forum) LoginWeb(w http.ResponseWriter, r *http.Request) {
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
		if err != nil {
			fmt.Printf("LoginWeb CreateSession error: %+v\n", err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken.String(),
			Expires: expiresAt,
			// MaxAge:  2 * int(time.Hour),
		})
		// fmt.Println(data.GetSession())
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
	err := tpl.ExecuteTemplate(w, "signup.html", nil)
	if err != nil {
		fmt.Printf("GetSignupPage Execute (signup.html) error: %+v\n", err)
	}
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
	err := r.ParseForm() // parses sign up form to fetch needed information
	if err != nil {
		fmt.Printf("SignUpUser ParseForm error: %+v\n", err)
	}

	var user User

	user.Email = r.FormValue("email")
	// check if e-mail is valid format
	isValidEmail := true

	if isValidEmail != strings.Contains(user.Email, "@") || isValidEmail != strings.Contains(user.Email, ".") { // checks if e-mail is valid by checking if it contains "@"
		isValidEmail = false
	}

	if !isValidEmail {
		fmt.Println("Email invalid")
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

	// fmt.Println(nameLength)

	// check pw criteria
	user.Password = r.FormValue("password")

	// fmt.Println(user)
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
	err2 := row.Scan(&username)
	if err2 != sql.ErrNoRows {
		fmt.Printf("sql scan row user error: %+v\n", err2)
		err3 := tpl.ExecuteTemplate(w, "signup.html", "username taken")
		if err3 != nil {
			fmt.Printf("SignUpUser ExecuteTemplate (username) error1: %+v\n", err3)
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
		err2 := tpl.ExecuteTemplate(w, "signup.html", "there was an error registering account")
		if err2 != nil {
			fmt.Printf("SignUpUser ExecuteTemplate (password) error:  %+v\n", err2)
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
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (data *Forum) GetUsernameFromSessionID(writer http.ResponseWriter, request *http.Request) string {
	c, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Printf("GetUsernameFromSessionID (Cookie) error: %+v\n", err)
			return ""
		}
	}

	sessionToken := c.Value
	a, err := data.GetSession()
	if err != nil {
		fmt.Printf("GetUsernameFromSessionID (GetSession) error: %+v\n", err)
		return ""
	}

	for _, sess := range a {
		// fmt.Println(sessionToken, " : ", sess.SessionID)
		if sessionToken == sess.SessionID {
			return sess.Username
		}
	}
	return ""
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
	a, err := data.GetSession()
	if err != nil {
		fmt.Printf("CheckCookie (GetSession) error: %+v\n", err)
	}

	for _, sess := range a {
		if sessionToken == sess.SessionID {
			return true
		}
	}
	return false
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
			_, err := data.DB.Exec("DELETE FROM session where sessionID ='" + currentSession.SessionID + "'")
			if err != nil {
				fmt.Printf("Logout Exec error: %+v\n", err)
			}
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
		data, err := data.GetPosts()
		if err != nil {
			fmt.Printf("HomePage (GetPosts) error: %+v\n", err)
		}
		err = tpl.ExecuteTemplate(writer, "guest.html", data)
		if err != nil {
			fmt.Printf("HomePage ExecuteTemplate (guest.html) error: %+v\n", err)
		}
		return
	} else {
		post, err := data.GetPosts()
		if err != nil {
			fmt.Printf("HomePage GetPosts (post) error: %+v\n", err)
		}
		var lastPost PostFeed
		if len(post) > 0 {
			lastPost = post[len(post)-1]
		}

		postCategory := request.FormValue("category")
		postTitle := request.FormValue("title")
		postContent := request.FormValue("content")

		postLikes := 0
		postDislikes := 0
		time := time.Now()
		postCreated := time.Format("01-02-2006 15:04")

		// checks session and selects the last one (the latest one)
		sess, err := data.GetSession()
		if err != nil {
			fmt.Printf("HomePage GetSession error: %+v\n", err)
		}
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
			postAndSession.Post, err = data.GetPosts()
			if err != nil {
				fmt.Printf("HomePage GetPosts (lastPost.Content) error: %+v\n", err)
			}
			err := tpl.ExecuteTemplate(writer, "./home", postAndSession)
			if err != nil {
				fmt.Printf("HomePage Execute (./home) error: %+v\n", err)
			}
			return
		} else {
			// postAndSession.UserSession = data.GetSession()[0]
			if postTitle != "" || postContent != "" || postCategory != "" {
				err := data.CreatePost(PostFeed{
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

				postAndSession.Post, err = data.GetPosts()
				if err != nil {
					fmt.Printf("HomePage (GetPosts) items error: %+v\n", err)
					return
				}

				err = tpl.ExecuteTemplate(writer, "home.html", postAndSession)
				if err != nil {
					fmt.Printf("HomePage ExecuteTemplate user homepage error: %+v\n", err)
					return
				}
				return
			}
		}
		data, err := data.GetPosts()
		postAndSession.Post = data
		if err != nil {
			fmt.Printf("HomePage (GetPosts) data error: %+v\n", err)
			return
		}
		err = tpl.ExecuteTemplate(writer, "home.html", postAndSession)
		if err != nil {
			fmt.Printf("HomePage ExecuteTemplate (home.html) error: %+v\n", err)
			return
		}
		return
	}
}

func (data *Forum) GuestView(writer http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	items, err := data.GetPosts()
	if err != nil {
		fmt.Printf("GuestView (GetPosts) items error: %+v\n", err)
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

	posts, err := data.GetPosts()
	if err != nil {
		fmt.Printf("CategoryDump (GetPosts) posts error: %+v\n", err)
		return
	} // get all posts
	// fmt.Println(posts)
	// check every post to find ones whose category matches our url path
	categoryFound := false // used to check if a valid category was entered
	for _, post := range posts {
		// fmt.Println(cat, post.Category)
		// fmt.Println(post.Category)
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
	err := tpl.ExecuteTemplate(writer, "passwordReset.html", nil)
	if err != nil {
		fmt.Printf("PwReset Execute (passwordReset.html) error: %+v\n", err)
	}
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

	post, err := data.GetPosts() // get all posts
	if err != nil {
		fmt.Printf("Threads (GetPosts) posts error: %+v\n", err)
		return
	}
	// TODO: ERROR HANDLING
	sess, _ := data.GetSession()
	currentSession := sess[len(sess)-1]
	// data.GetSession()[len(data.GetSession())-1]
	cmnt, _ := data.GetComments()
	var lastComment Comment
	if len(cmnt) > 0 {
		lastComment = cmnt[len(cmnt)-1]
	}

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
	post, err := data.GetPosts() // get all posts
	if err != nil {
		fmt.Printf("ThreadGuest (GetPosts) posts error: %+v\n\n", err)
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

	user, err := data.GetSession()
	if err != nil {
		fmt.Printf("UserPosts GetSession error: %+v\n", err)
	}

	currentUser := user[len(user)-1]
	// if user.session == user in post --- send this post

	posts, err := data.GetPosts()
	if err != nil {
		fmt.Printf("UserPosts GetPosts error: %+v\n", err)
		return
	}

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
	err = tpl.ExecuteTemplate(writer, "posts.html", usrPosts)
	if err != nil {
		fmt.Printf("UserPosts ExecuteTemplate (posts.html) error: %+v\n", err)
		return
	}
}

func (data *Forum) UserLikes(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	user, err := data.GetSession()
	if err != nil {
		fmt.Printf("UserLikes GetSession error: %+v\n", err)
	}

	currentUser := user[len(user)-1]

	allReactions, err := data.GetReactions()
	if err != nil {
		fmt.Printf("UserLikes (GetReactions) error: %+v\n", err)
		return
	}

	posts, err := data.GetPosts()
	if err != nil {
		fmt.Printf("UserLikes GetPosts error: %+v\n", err)
		return
	}

	comments, err := data.GetComments()
	if err != nil {
		fmt.Printf("UserLikes GetComments error: %+v\n", err)
		return
	}

	type UserLikes struct {
		Posts    []PostFeed
		Comments []Comment
	}
	var likesByUser UserLikes

	for _, post := range posts {
		if post.Username == currentUser.Username {
			for _, reaction := range allReactions {
				if reaction.Liked && post.Username == reaction.Username && post.PostID == reaction.PostID {
					likesByUser.Posts = append(likesByUser.Posts, post)
				}
			}
		}
	}

	for _, comment := range comments {
		if comment.UserId == currentUser.Username {
			for _, reaction := range allReactions {
				if reaction.Liked && comment.UserId == reaction.Username && comment.CommentID == reaction.CommentID {
					likesByUser.Comments = append(likesByUser.Comments, comment)
				}
			}
		}
	}
	err = tpl.ExecuteTemplate(writer, "likes.html", likesByUser)
	if err != nil {
		fmt.Printf("UserLikes ExecuteTemplate (likes.html) error: %+v\n", err)
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

func (data *Forum) HandleLikeDislikeForPost(writer http.ResponseWriter, request *http.Request, isLike bool) {
	loggedIn := data.CheckCookie(writer, request)
	if !loggedIn {
		fmt.Printf("Guests are unable to like/dislike posts\n")
		return
	}
	items, err := data.GetPosts()
	// fmt.Printf("Items: %+v\n", items)
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
		username := data.GetUsernameFromSessionID(writer, request)
		reaction, error := data.GetReactionByPostID(strconv.Itoa(requestedItem.PostID), username)
		if error != nil {
			fmt.Printf("HandleLikeDislike (GetReactionByPostID) error: %v\n", error)
		}
		if reaction == nil {
			err := data.CreateReaction(Reaction{
				PostID:   requestedItem.PostID,
				Username: username,
				Liked:    isLike,
				Disliked: !isLike,
			})
			if err != nil {
				fmt.Printf("HandleLikeDislike (CreateReaction) error: %v\n", err)
			}
		}
		if reaction != nil {
			err := data.UpdateReaction(Reaction{
				ReactionID: reaction.ReactionID,
				PostID:     requestedItem.PostID,
				Username:   username,
				Liked:      isLike,
				Disliked:   !isLike,
			})
			if err != nil {
				fmt.Printf("HandleLikeDislike (UpdateReaction) error: %v\n", err)
			}
		}

		// we need to know if anything actually happened
		// we check to see if the reaction will have changed based on what it was already compared with isLike
		success := false
		// if wasLikedBefore != clientAskedForLike OR wasDislikedBefore != clientAskedForDislike THEN success=true
		if reaction == nil {
			success = true
		} else if reaction.Liked != isLike || reaction.Disliked != !isLike {
			success = true
		}

		updatedLikes, err := getLikesForPost(data.DB, requestedItem.PostID)
		if err != nil {
			fmt.Printf("HandleLikeDislike (getLikesForPost) error: %v\n", err)
		}

		updatedDislikes, err := getDislikesForPost(data.DB, requestedItem.PostID)
		if err != nil {
			fmt.Printf("HandleLikeDislike (getLikesForPost) error: %v\n", err)
		}

		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write([]byte(fmt.Sprintf("{\"success\":%t,\"likes\":%d,\"dislikes\":%d}", success, updatedLikes, updatedDislikes)))
		if err != nil {
			fmt.Printf("unable to send json response for post %d\n", reqItemID)
			return
		}
		fmt.Printf("modified dis/likes on post %d\n", reqItemID)
	}
}

func (data *Forum) HandleLikeDislikeForComment(writer http.ResponseWriter, request *http.Request, isLike bool) {
	loggedIn := data.CheckCookie(writer, request)
	if !loggedIn {
		fmt.Printf("Guests are unable to like/dislike comments\n")
		return
	}
	items, err := data.GetComments()
	// fmt.Printf("Items: %+v\n", items)
	if err != nil {
		fmt.Printf("HandleLikeDislike (GetComments) comments error: %+v\n", err)
		return
	}

	reqItemIDraw := request.URL.Query().Get("id")
	reqItemID, err := strconv.Atoi(reqItemIDraw)
	if err != nil {
		fmt.Printf("unable to parse comment id: %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte("{\"500\": \"Error parsing comment id\"}"))
		if err != nil {
			fmt.Printf("unable to send json response for comments %d\n", reqItemID)
			return
		}
		return
	}
	requestedItem := Comment{}

	for _, item := range items {
		if item.CommentID == reqItemID {
			requestedItem = item
		}
	}

	if requestedItem.CreatedAt == "" {
		fmt.Printf("unable to find comment %d in db: %v\n", reqItemID, err)
		writer.WriteHeader(http.StatusNotFound)
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte("{\"404\": \"Error finding comment\"}"))
		if err != nil {
			fmt.Printf("unable to send json response for comment %d\n", reqItemID)
			return
		}
		return
	}

	j, err := requestedItem.MarshallJSON()
	if err != nil {
		fmt.Printf("unable to marshal json for comment %d: %v\n", reqItemID, err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte("{\"500\": \"Error marshalling json for comment\"}"))
		if err != nil {
			fmt.Printf("unable to send json response for comment %d\n", reqItemID)
			return
		}
		return
	}

	switch request.Method {
	case http.MethodGet:
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write(j)
		if err != nil {
			fmt.Printf("unable to send json response for comment %d\n", reqItemID)
		}
	case http.MethodPost:
		username := data.GetUsernameFromSessionID(writer, request)
		reaction, error := data.GetReactionByCommentID(strconv.Itoa(requestedItem.CommentID), username)
		if error != nil {
			fmt.Printf("HandleLikeDislike (GetReactionByCommentID) error: %v\n", error)
		}
		if reaction == nil {
			err := data.CreateReaction(Reaction{
				CommentID: requestedItem.CommentID,
				Username:  username,
				Liked:     isLike,
				Disliked:  !isLike,
			})
			if err != nil {
				fmt.Printf("HandleLikeDislike (CreateReaction) comment error: %v\n", err)
			}
		}
		if reaction != nil {
			err := data.UpdateReaction(Reaction{
				ReactionID: reaction.ReactionID,
				CommentID:  requestedItem.CommentID,
				Username:   username,
				Liked:      isLike,
				Disliked:   !isLike,
			})
			if err != nil {
				fmt.Printf("HandleLikeDislike (UpdateReaction) comment error: %v\n", err)
			}
		}
		success := false
		// if wasLikedBefore != clientAskedForLike OR wasDislikedBefore != clientAskedForDislike THEN success=true
		if reaction == nil {
			success = true
		} else if reaction.Liked != isLike || reaction.Disliked != !isLike {
			success = true
		}

		updatedLikes, err := getLikesForComment(data.DB, requestedItem.CommentID)
		if err != nil {
			fmt.Printf("HandleLikeDislike (getLikesForComment) error: %v\n", err)
		}

		updatedDislikes, err := getDislikesForComment(data.DB, requestedItem.CommentID)
		if err != nil {
			fmt.Printf("HandleLikeDislike (getLikesForComment) error: %v\n", err)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write([]byte(fmt.Sprintf("{\"success\":%t,\"likes\":%d,\"dislikes\":%d}", success, updatedLikes, updatedDislikes)))
		if err != nil {
			fmt.Printf("unable to send json response for comment %d\n", reqItemID)
			return
		}
		fmt.Printf("modified dis/likes on comment %d\n", reqItemID)
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
	// case "/comments":
	// data.UserComments(w, r)
	case "/likes":
		data.UserLikes(w, r)
	// case "/shares":
	// data.UserShares(w, r)
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
	case "/likePost":
		data.HandleLikeDislikeForPost(w, r, true)
	case "/dislikePost":
		data.HandleLikeDislikeForPost(w, r, false)
	case "/likeComment":
		data.HandleLikeDislikeForComment(w, r, true)
	case "/dislikeComment":
		data.HandleLikeDislikeForComment(w, r, false)
	}
}
