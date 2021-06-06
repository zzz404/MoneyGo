package deposit

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/zzz404/MoneyGo/internal/db"
)

type CoinType string

const (
	RMB CoinType = "RMB" // 人民幣
	TWD CoinType = "TWD" // 台幣
	USD CoinType = "USD" // 美金
)

type DepositType int8

const (
	DemandDeposit DepositType = 1 // 活存
	TimeDeposit   DepositType = 2 // 定存
)

type Deposit struct {
	Id       int
	MemberId int8
	BankId   int8
	Type     DepositType
	Amount   float32
	CoinType CoinType
}

type Dao struct {
	Columns []string
}

var columns = []string{"id", "memberId", "bankId", "type", "amount", "coinType"}

func (d *Deposit) load(rows *sql.Rows) error {
	return rows.Scan(&d.Id, &d.MemberId, &d.BankId, &d.Type, &d.Amount, &d.CoinType)
}

func (d *Deposit) toValues() []interface{} {
	return []interface{}{d.Id, d.MemberId, d.BankId, d.Type, d.Amount, d.CoinType}
}

func QueryDeposits(memberId int) ([]*Deposit, error) {
	sql := "SELECT " + db.ToColumnsString(columns) + " FROM Deposit WHERE memberId=?"
	rows, err := db.DB.Query(sql, memberId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deposits []*Deposit
	for rows.Next() {
		deposit := &Deposit{}
		err = deposit.load(rows)
		if err != nil {
			return nil, err
		}
		deposits = append(deposits, deposit)
	}
	return deposits, nil

}

func GetDeposit(id int) (*Deposit, error) {
	sql := "SELECT " + db.ToColumnsString(columns) + " FROM Deposit WHERE id=?"
	rows, err := db.DB.Query(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deposit *Deposit = nil
	for rows.Next() {
		if deposit == nil {
			deposit = &Deposit{}
		} else {
			return nil, errors.New(fmt.Sprintf("Deposit id %d 不只一個!?", id))
		}
		err = deposit.load(rows)
		if err == nil {
			return nil, err
		}
	}
	return deposit, nil
}

func AddDeposit(deposit *Deposit) error {
	params, err := db.ToSqlParams(len(columns) - 1)
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("INSERT INTO Deposit (%s) VALUES (%s)",
		db.ToColumnsString(columns[1:]), params)
	pstmt, err := db.DB.Prepare(sql)
	defer pstmt.Close()

	_, err = pstmt.Exec(deposit.toValues()[1:]...)
	return err
}

func UpdateDeposit(deposit *Deposit) error {
	sql := fmt.Sprintf("UPDATE Deposit SET %s WHERE id=?",
		db.ToSettersString(columns[1:]))
	pstmt, err := db.DB.Prepare(sql)
	defer pstmt.Close()

	values := deposit.toValues()[1:]
	values = append(values, deposit.Id)
	_, err = pstmt.Exec(values...)
	return err
}

func DeleteDeposit(id int) error {
	sql := "DELETE FROM Deposit WHERE id=?"
	pstmt, err := db.DB.Prepare(sql)
	defer pstmt.Close()

	_, err = pstmt.Exec(id)
	return err
}
