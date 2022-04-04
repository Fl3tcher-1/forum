package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var login *template.Template
var home *template.Template
var categories *template.Template
var reset *template.Template
var signup *template.Template
var profile *template.Template

var photo *template.Template
var posts *template.Template
var comments *template.Template
var likes *template.Template
var shares *template.Template
var userinfo *template.Template
var custom *template.Template

func main() {

	login = template.Must(template.ParseFiles("./templates/login.html"))
	home = template.Must(template.ParseFiles("./templates/home.html"))
	categories = template.Must(template.ParseFiles("./templates/categories.html"))
	reset = template.Must(template.ParseFiles("./templates/passwordReset.html"))
	signup = template.Must(template.ParseFiles("./templates/signup.html"))
	profile = template.Must(template.ParseFiles("./templates/profile.html"))

	photo = template.Must(template.ParseFiles("./templates/photo.html")) // can these be used in a struct format???
	posts = template.Must(template.ParseFiles("./templates/posts.html"))
	comments = template.Must(template.ParseFiles("./templates/comments.html"))
	likes = template.Must(template.ParseFiles("./templates/likes.html"))
	shares = template.Must(template.ParseFiles("./templates/shares.html"))
	userinfo = template.Must(template.ParseFiles("./templates/userinfo.html"))
	custom = template.Must(template.ParseFiles("./templates/customize.html"))

	mux := http.NewServeMux()

	mux.HandleFunc("/login", loginWeb)
	mux.HandleFunc("/home", homePage)
	mux.HandleFunc("/categories", categoriesList)
	mux.HandleFunc("/reset", pwReset)
	mux.HandleFunc("/signup", signUp)
	mux.HandleFunc("/profile", userProfile)

	mux.HandleFunc("/photo", userPhoto)
	mux.HandleFunc("/posts", userPosts)
	mux.HandleFunc("/comments", userComments)
	mux.HandleFunc("/likes", userLikes)
	mux.HandleFunc("/shares", userShares)
	mux.HandleFunc("/info", userInfo)
	mux.HandleFunc("/custom", customization)

	fmt.Printf("Starting server at port 8080\n\t -----------\nhttp://localhost:8080/login\n")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(500, "500 Internal server error", err)
	}

	fmt.Println("hi")
}

func loginWeb(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	login.Execute(writer, nil)

}
func homePage(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	home.Execute(writer, nil)
}
func categoriesList(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	categories.Execute(writer, nil)
}
func pwReset(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	reset.Execute(writer, nil)
}
func signUp(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	signup.Execute(writer, nil)
}
func userProfile(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	profile.Execute(writer, nil)
}
func userPhoto(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	photo.Execute(writer, nil)
}
func userPosts(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	posts.Execute(writer, nil)
}
func userComments(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	comments.Execute(writer, nil)
}
func userLikes(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	likes.Execute(writer, nil)
}
func userShares(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	shares.Execute(writer, nil)
}
func userInfo(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	userinfo.Execute(writer, nil)
}
func customization(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	custom.Execute(writer, nil)
}
