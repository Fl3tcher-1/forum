package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

var DB *sql.DB

func UserDatabase() {
	db, _ := sql.Open("sqlite3", "userdata.db")
	// if err != nil {
	// 	checkErr(err)
	// }

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT, email TEXT)")
	statement.Exec()

	// if err != nil {
	// 	checkErr(err)
	// }

	DB = db
}

func checkErr(err error) {
	log.Fatal(err)
}
