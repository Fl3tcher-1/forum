package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/database"

	// "forum/database"

	_ "github.com/mattn/go-sqlite3"
)

func cssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/stylesheet.css")
}

func main() {
	db, err := sql.Open("sqlite3", "./database/userdata.db")
	if err != nil {
		fmt.Printf("main (sql.Open) error: %+v\n", err)
		os.Exit(1)
	}
	data := database.Connect(db)

	mux := http.NewServeMux()
	mux.Handle("/pics/", http.StripPrefix("/pics", http.FileServer(http.Dir("./pics"))))
	mux.HandleFunc("/", data.Handler)
	mux.HandleFunc("/category/", data.CategoryDump)
	mux.HandleFunc("/categoryg/", data.CategoryDump)
	mux.HandleFunc("/threadg/", data.ThreadGuest)
	mux.HandleFunc("/thread/", data.Threads)
	mux.HandleFunc("/category/stylesheet", cssHandler)
	mux.HandleFunc("/threadg/stylesheet", cssHandler)
	mux.HandleFunc("/thread/stylesheet", cssHandler)
	mux.HandleFunc("/categoryg/stylesheet", cssHandler)

	fmt.Println("Starting server at port 8080:\n http://localhost:8080/login")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(500, "500 Internal server error:", err)
		fmt.Printf("main ListenAndServe error: %+v\n", err)
	}
}
