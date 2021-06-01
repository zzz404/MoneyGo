package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zzz404/MoneyGo/internal/db"
)

func main() {
	members := db.QueryMembers()
	for _, m := range members {
		fmt.Println(m)
	}
}
