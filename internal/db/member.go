package db

import (
	"database/sql"
)

type Member struct {
	Id   int8
	Name string
}

func QueryMembers() []Member {
	db, err := sql.Open("sqlite3", "D:/My Data/Money/MoneyGo.db")
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM Member")
	checkErr(err)

	var members []Member
	for rows.Next() {
		member := Member{}
		err = rows.Scan(&member.Id, &member.Name)
		checkErr(err)
		members = append(members, member)
	}
	return members
}
