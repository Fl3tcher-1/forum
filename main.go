package main

import (
	"database/sql"
	"fmt"
	"forum/database"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// handles possible web directories
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	// page handlers
	case "/stylesheet": // handle css
		http.ServeFile(w, r, "./templates/stylesheet.css")
	case "/":
		database.LoginWeb(w, r)
	case "/login":
		database.LoginWeb(w, r)
	case "/home":
		database.HomePage(w, r)
	case "/categories":
		database.CategoriesList(w, r)
	case "/reset":
		database.PwReset(w, r)
	case "/signup":
		database.GetSignupPage(w, r)
	case "/sign-up-form":
		database.SignUpUser(w, r)
	case "/profile":
		database.UserProfile(w, r)
	case "/thread":
		database.Threads(w, r)
	case "/about":
		database.AboutFunc(w, r)
	case "/contact-us":
		database.ContactUs(w, r)

		// user handlers
	case "/photo":
		database.UserPhoto(w, r)
	case "/posts":
		database.UserPosts(w, r)
	case "/comments":
		database.UserComments(w, r)
	case "/likes":
		database.UserLikes(w, r)
	case "/dislikes":
		database.UserDislikes(w, r)
	case "/shares":
		database.UserShares(w, r)
	case "/info":
		database.UserInfo(w, r)
	case "/custom":
		database.Customization(w, r)

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
		database.AddLike(w, r)
	case "/dislike":
		database.AddDislike(w, r)
	}
}

func main() {
	database.UserDatabase()

	// posts, err := sql.Open("sqlite3", "./database/feed.db")
	// if err != nil {
	// 	fmt.Printf("main Database sql.Open error: %+v\n", err)
	// }
	// feed := database.Feed(posts)

	// feed.Add(database.PostFeed{
	// 	Content: "the monkeys are taking control",
	// })

	// items := feed.Get()
	// fmt.Println(items)

	commentDB, err := sql.Open("sqlite3", "./database/comments.db")
	if err != nil {
		fmt.Printf("comments Database sql.Open error: %+v\n", err)
	}

	comments := database.Comments(commentDB)

	// comments.AddComment(database.PostFeed{
	// 	Title: "monke",
	// 	Content: "monkeys are taking control",
	// 	Likes: 3,
	// 	Created: "now",
	// 	Category: "",
	// })

	c := comments.GetComments()
	fmt.Println(c)
	fmt.Println(comments)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	fmt.Println("Starting server at port 8080: http://localhost:8080/login")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(500, "500 Internal server error:", err)
		fmt.Printf("main ListenAndServe error: %+v\n", err)
	}
}
