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
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", data.Handler)

	fmt.Println("Starting server at port 8080: http://localhost:8080/login")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(500, "500 Internal server error:", err)
		fmt.Printf("main ListenAndServe error: %+v\n", err)
	}

	fmt.Println("hi")
}
