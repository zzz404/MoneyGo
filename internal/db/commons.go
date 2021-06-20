package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
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

type SqlBuilder struct {
	Columns    []string
	Tables     []string
	Conditions []string
	Variables  []interface{}
	OrderBys   []string
}

func (sb *SqlBuilder) AddTable(table string) *SqlBuilder {
	sb.Tables = append(sb.Tables, table)
	return sb
}

func (sb *SqlBuilder) SetColumns(columns []string) *SqlBuilder {
	sb.Columns = columns
	return sb
}

func (sb *SqlBuilder) AddCondition(cond string, variable interface{}) *SqlBuilder {
	sb.Conditions = append(sb.Conditions, cond)
	sb.Variables = append(sb.Variables, variable)
	return sb
}

func (sb *SqlBuilder) AddOrderBy(orderby string) *SqlBuilder {
	sb.OrderBys = append(sb.OrderBys, orderby)
	return sb
}

func (sb *SqlBuilder) BuildSql() string {
	sql := fmt.Sprintf("SELECT %s FROM %s", strings.Join(sb.Columns, ", "), strings.Join(sb.Tables, ", "))
	if len(sb.Conditions) > 0 {
		sql += " WHERE " + strings.Join(sb.Conditions, " AND ")
	}
	if len(sb.OrderBys) > 0 {
		sql += " ORDER BY " + strings.Join(sb.OrderBys, ", ")
	}
	return sql
}
