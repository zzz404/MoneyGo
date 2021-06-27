package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/zzz404/MoneyGo/internal/utils"
)

var DB *sql.DB

func init() {
	bt, err := ioutil.ReadFile("cfg/db-path.txt")
	if err != nil {
		panic(err)
	}
	dbPath := string(bt)

	if err := copy_bak(dbPath); err != nil {
		panic(err)
	}

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
}

func copy_bak(dbPath string) error {
	dir := dbPath + ".bak"
	if err := utils.AssertDirExists(dir); err != nil {
		return err
	}

	timestamp := utils.GetTimeMillis()
	bakFileName := fmt.Sprintf("%s.%d", filepath.Base(dbPath), timestamp)
	bakPath := fmt.Sprintf("%s/%s", dir, bakFileName)

	return utils.CopyFile(dbPath, bakPath)
}

func ToSqlParams(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("n 必須大於 0")
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

type QueryAble interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	// Query(query string, args ...interface{}) (*sql.Rows, error)
	// QueryRow(query string, args ...interface{}) *sql.Row
}

func ExecuteSql(sql string, args ...interface{}) (sql.Result, error) {
	pstmt, err := DB.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer pstmt.Close()

	return pstmt.Exec(args...)
}

type SqlExecuter struct {
	tx *sql.Tx
}

func (exe *SqlExecuter) StartTx() error {
	var err error
	if exe.tx != nil {
		err = exe.tx.Rollback()
		if err != nil {
			return err
		}
	}
	exe.tx, err = DB.Begin()
	return err
}

func (exe *SqlExecuter) Commit() error {
	if exe.tx == nil {
		return nil
	}
	err := exe.tx.Commit()
	if err == nil {
		exe.tx = nil
	}
	return err
}

func (exe *SqlExecuter) Rollback() error {
	if exe.tx == nil {
		return nil
	}
	err := exe.tx.Rollback()
	if err == nil {
		exe.tx = nil
	}
	return err
}

func (exe *SqlExecuter) Close() error {
	if exe.tx != nil {
		return exe.tx.Rollback()
	}
	return nil
}

func (exe *SqlExecuter) ExecuteSql(sql_of_update string, args ...interface{}) (sql.Result, error) {
	var pstmt *sql.Stmt
	var err error

	if exe.tx == nil {
		pstmt, err = DB.Prepare(sql_of_update)
	} else {
		pstmt, err = exe.tx.Prepare(sql_of_update)
	}
	if err != nil {
		return nil, err
	}
	defer pstmt.Close()

	return pstmt.Exec(args...)
}
