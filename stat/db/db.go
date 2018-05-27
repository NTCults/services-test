package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Connect *sql.DB

func ConnectToDB() {
	var err error
	Connect, err = sql.Open("sqlite3", "./db/stat.db")
	if err != nil {
		fmt.Println(err)
	}
	err = Connect.Ping()
	if err != nil {
		log.Fatal(err)
	}
	statement, _ := Connect.Prepare(`CREATE TABLE IF NOT EXISTS stat (
		account TEXT, campaign_id INTEGER PRIMARY KEY, shows INTEGER, clicks INTEGER, costs INTEGER, date TEXT)`)
	statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
