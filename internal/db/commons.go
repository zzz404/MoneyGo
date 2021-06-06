package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	dbPath, err := ioutil.ReadFile("cfg/db-path.txt")
	assertSucc(err)
	db1, err := sql.Open("sqlite3", string(dbPath))
	assertSucc(err)
	DB = db1
}

func assertSucc(err error) {
	if err != nil {
		panic(err)
	}
}

func ToSqlParams(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("n 必須大於 0!")
	} else if n == 1 {
		return "?", nil
	} else {
		return strings.Repeat("?, ", n-1) + "?", nil
	}
}

func ToColumnsString(columns []string) string {
	return strings.Join(columns, ", ")
}

func ToSettersString(columns []string) string {
	switch len(columns) {
	case 0:
		return ""
	case 1:
		return columns[0] + "=?"
	}
	builder := strings.Builder{}
	fmt.Fprint(&builder, columns[0]+"=?")
	for _, column := range columns[1:] {
		fmt.Fprintf(&builder, ", %s=?", column)
	}
	return builder.String()
}
