package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	dbPath, err := ioutil.ReadFile("cfg/db-path.txt")
	assertSucc(err)
	db1, err := sql.Open("sqlite3", string(dbPath))
	assertSucc(err)
	db = db1
}

func close() {
	fmt.Println("---------closed")
	db.Close()
}

func assertSucc(err error) {
	if err != nil {
		panic(err)
	}
}
