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
	Connect, err = sql.Open("sqlite3", "./db/campaigns.db")
	if err != nil {
		fmt.Println(err)
	}
	err = Connect.Ping()
	if err != nil {
		log.Fatal(err)
	}
	statement, _ := Connect.Prepare("CREATE TABLE IF NOT EXISTS campaigns (account TEXT, id INTEGER PRIMARY KEY, title TEXT)")
	statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
