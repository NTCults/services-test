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
	Connect, err = sql.Open("sqlite3", "./db/tags.db")
	if err != nil {
		fmt.Println(err)
	}
	err = Connect.Ping()
	if err != nil {
		log.Fatal(err)
	}
	statement, _ := Connect.Prepare(`CREATE TABLE IF NOT EXISTS tags (
		account TEXT, campaign_id INTEGER, tag TEXT)`)
	statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
