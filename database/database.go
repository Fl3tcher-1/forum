package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

var DB *sql.DB

func UserDatabase() {
	db, err := sql.Open("sqlite3", "./database/userdata.db")
	if err != nil {
		checkErr(err)
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT, email TEXT, passwordHASH TEXT)")
	statement.Exec()

	if err != nil {
		checkErr(err)
	}

	DB = db
}

func checkErr(err error) {
	log.Fatal(err)
}

// func Feed() {
// 	db, err := sql.Open("sqlite3", "./database/feed.db")
// 	if err != nil {
// 		checkErr(err)
// 	}

// 	stmt, err := db.Prepare(`CREATE TABLE "feed" (
// 	"iD"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
// 	"title"	TEXT,
// 	"content"	TEXT,
// 	"comments"	TEXT,
// 	"likes"	INTEGER,
// 	"created"	TEXT,
// 	"category"	TEXT
// );`)
// 	stmt.Exec()
// 	if err != nil {
// 		checkErr(err)
// 	}
// }
