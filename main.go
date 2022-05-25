package main

import (
	"database/sql"
	"fmt"
	"forum/database"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./database/userdata.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	data := database.Connect(db)

	// posts, err := sql.Open("sqlite3", "./database/feed.db")
	// if err != nil {
	// 	fmt.Printf("main Database sql.Open error: %+v\n", err)
	// }
	// feed := database.Feed(posts)

	// feed.Add(database.PostFeed{
	// 	Content: "the monkeys are taking control",
	// })

	// items :=feed.Get()
	// fmt.Println(items)

	// commentDB, err := sql.Open("sqlite3", "./database/comments.db")
	// if err != nil {
	// 	fmt.Printf("comments Database sql.Open error: %+v\n", err)
	// }

	// comments := database.Comments(commentDB)

	// comments.AddComment(database.PostFeed{
	// 	Title: "monke",
	// 	Content: "monkeys are taking control",
	// 	Likes: 3,
	// 	Created: "now",
	// 	Category: "",
	// })

	// c := comments.GetComments()
	// fmt.Println(c)

	// fmt.Println(comments)

	mux := http.NewServeMux()
	mux.HandleFunc("/", data.Handler)

	fmt.Println("Starting server at port 8080: http://localhost:8080/login")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(500, "500 Internal server error:", err)
		fmt.Printf("main ListenAndServe error: %+v\n", err)
	}

	fmt.Println("hi")
}
