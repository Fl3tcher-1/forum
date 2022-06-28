package sesman

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserName string
	Password []byte
	Email    string
}

type session struct {
	un           string
	lastActivity time.Time
}

var (
	dbUsers           = map[string]User{}    // user ID, user
	dbSessions        = map[string]session{} // session ID, session
	dbSessionsCleaned time.Time
)

const sessionLength int = 30

func init() {
	_, err := template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Printf("init (ParseGlob) error: %+v\n", err)
		return
	}
	dbSessionsCleaned = time.Now()
	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("init ListenAndServe error: %+v\n", err)
	}
}

func getUser(w http.ResponseWriter, req *http.Request) User {
	// get cookie
	c, err := req.Cookie("session")
	if err != nil {
		sID := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		fmt.Printf("getUser (Cookie) error: %+v\n", err)
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	var u User
	// if the user exists already, get user

	if s, ok := dbSessions[c.Value]; ok {
		s.lastActivity = time.Now()
		dbSessions[c.Value] = s
		u = dbUsers[s.un]
	}
	return u
}

func alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		fmt.Printf("alreadyLoggedIn (Cookie) error: %+v\n", err)
		return false
	}
	s, ok := dbSessions[c.Value]
	if ok {
		s.lastActivity = time.Now()
		dbSessions[c.Value] = s
	}
	_, ok = dbUsers[s.un]
	// refresh session
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	return ok
}

func cleanSessions() {
	fmt.Println("BEFORE CLEAN") // for demonstration purposes
	showSessions()              // for demonstration purposes
	for k, v := range dbSessions {
		if time.Since(v.lastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()
	fmt.Println("AFTER CLEAN") // for demonstration purposes
	showSessions()             // for demonstration purposes
}

// for demonstration purposes
func showSessions() {
	fmt.Println("********")
	for k, v := range dbSessions {
		fmt.Println(k, v.un)
	}
	fmt.Println("")
}

func index(w http.ResponseWriter, req *http.Request) {
	// U := getUser(w, req)
	var u User
	tpl, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		log.Println(err.Error(), u, "")
		return
	}
	c, err := req.Cookie("session")
	if err != nil {
		sID := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		fmt.Printf("index (Cookie) error: %+v\n", err)
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	showSessions() // for demonstration purposes
	err = tpl.Execute(w, u)
	if err != nil {
		log.Println(err.Error(), "")
		return
	}
}

func signup(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	var u User
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("email")

		err := req.ParseForm()
		if err != nil {
			fmt.Printf("signup ParseForm error: %+v\n", err)
		}
		// username taken?
		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
		// create session
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un, time.Now()}
		// store user in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		u = User{un, bs, f}
		dbUsers[un] = u
		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	showSessions() // for demonstration purposes
	tpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		log.Println(err.Error(), u, "")
		return
	}
	err = tpl.Execute(w, u)
	if err != nil {
		log.Println(err.Error(), "")
		return
	}
}

func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	var u User
	// process form submission
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		err := req.ParseForm()
		if err != nil {
			fmt.Printf("login ParseForm error: %+v\n", err)
		}
		// is there a username?
		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// does the entered password match the stored password?
		err2 := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err2 != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// create session
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un, time.Now()}
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	showSessions() // for demonstration purposes
	tpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Println(err.Error(), "")
		return
	}
	err = tpl.Execute(w, u)
	if err != nil {
		log.Println(err.Error(), "")
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	c, err := r.Cookie("session")
	if err != nil {
		fmt.Printf("logout (Cookie) error: %+v\n", err)
	}
	// delete the session
	delete(dbSessions, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// clean up dbSessions
	if time.Since(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
