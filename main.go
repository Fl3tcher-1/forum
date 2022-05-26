package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"v2/Forum/database"
	"v2/Forum/endpoints"

	_ "github.com/mattn/go-sqlite3"
)

// handles possible endpoints
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	//page handlers
	case "/stylesheet": //handle css
		http.ServeFile(w, r, "./templates/stylesheet.css")
	case "/":
		endpoints.LoginWeb(w, r)
	case "/login":
		endpoints.LoginWeb(w, r)
	case "/home":
		endpoints.HomePage(w, r)
	case "/categories":
		endpoints.CategoriesList(w, r)
	case "/reset":
		endpoints.PwReset(w, r)
	case "/signup":
		endpoints.GetSignupPage(w, r)
	case "/sign-up-form":
		endpoints.SignUpUser(w, r)
	case "/profile":
		endpoints.UserProfile(w, r)
	case "/thread":
		endpoints.Threads(w, r)
	case "/about":
		endpoints.AboutFunc(w, r)
	case "/contact-us":
		endpoints.ContactUs(w, r)

		// user handlers
	case "/photo":
		endpoints.UserPhoto(w, r)
	case "/posts":
		endpoints.UserPosts(w, r)
	case "/comments":
		endpoints.UserComments(w, r)
	case "/likes":
		endpoints.UserLikes(w, r)
	case "/shares":
		endpoints.UserShares(w, r)
	case "/info":
		endpoints.UserInfo(w, r)
	case "/custom":
		endpoints.Customization(w, r)

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

func main() {

	database.UserDatabase()

	// posts, err := sql.Open("sqlite3", "./database/feed.db")
	// if err != nil {
	// 	database.CheckErr(err)
	// }
	// feed := database.Feed(posts)

	// feed.Add(database.PostFeed{
	// 	Content: "the monkeys are taking control",
	// })

	// items :=feed.Get()
	// fmt.Println(items)

	commentDB, _ :=sql.Open("sqlite3", "./database/comments.db")
	
	comments := database.Comments(commentDB)
	
	// comments.AddComment(database.PostFeed{
	// 	Title: "monke",
	// 	Content: "monkeys are taking control",
	// 	Likes: 3,
	// 	Created: "now",
	// 	Category: "",
	// })

	c:=comments.GetComments()
	fmt.Println(c)
	
	// fmt.Println(comments)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	fmt.Printf("Starting server at port 8080\n\t -----------\nhttp://localhost:8080/login\n")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(500, "500 Internal server error", err)
	}

}
