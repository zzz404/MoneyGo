package deposit

import (
	"database/sql"
	"fmt"

	bk "github.com/zzz404/MoneyGo/internal/bank"
	"github.com/zzz404/MoneyGo/internal/db"
	mb "github.com/zzz404/MoneyGo/internal/member"
	"github.com/zzz404/MoneyGo/internal/utils"
)

type DepositType struct {
	Code int
	Name string
}

var DemandDeposit DepositType = DepositType{Code: 1, Name: "活存"}
var TimeDeposit DepositType = DepositType{Code: 2, Name: "定存"}

var DepositTypes = []*DepositType{&DemandDeposit, &TimeDeposit}

func GetDepositTypeByCode(code int) (*DepositType, error) {
	switch code {
	case DemandDeposit.Code:
		return &DemandDeposit, nil
	case TimeDeposit.Code:
		return &TimeDeposit, nil
	}
	return nil, fmt.Errorf("不認識的 DepositType %d", code)
}

type Deposit struct {
	Id       int
	Member   *mb.Member
	BankId   int
	Type     *DepositType
	Amount   float32
	CoinType *utils.CoinType
}

func (d *Deposit) BankName() string {
	bank := bk.GetBank(d.BankId)
	if bank == nil {
		return ""
	} else {
		return bank.Name
	}
}

type Dao struct {
	Columns []string
}

var columns = []string{"id", "memberId", "bankId", "type", "amount", "coinType"}

func (d *Deposit) load(rows *sql.Rows) error {
	var coinTypeCode string
	var depositTypeCode int

	var memberId int
	err := rows.Scan(&d.Id, &memberId, &d.BankId, &depositTypeCode, &d.Amount, &coinTypeCode)
	if err != nil {
		return err
	}

	member, err := mb.GetMember(memberId)
	if err != nil {
		return err
	}
	d.Member = member

	d.Type, err = GetDepositTypeByCode(depositTypeCode)
	if err != nil {
		return err
	}
	d.CoinType, err = utils.GetCoinTypeByCode(coinTypeCode)
	return err
}

func (d *Deposit) toTableValues() []interface{} {
	return []interface{}{d.Id, d.Member.Id, d.BankId, d.Type.Code, d.Amount, d.CoinType.Code}
}

func QueryDeposits(memberId int) ([]*Deposit, error) {
	sql := "SELECT " + db.ToColumnsString(columns) + " FROM Deposit WHERE memberId=? ORDER BY bankId"
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
			return nil, fmt.Errorf("Deposit id %d 不只一個!?", id)
		}
		err = deposit.load(rows)
		if err != nil {
			return nil, err
		}
	}
	return deposit, nil
}

func AddDeposit(deposit *Deposit) (int, error) {
	params, err := db.ToSqlParams(len(columns) - 1)
	if err != nil {
		return 0, err
	}
	sql := fmt.Sprintf("INSERT INTO Deposit (%s) VALUES (%s)",
		db.ToColumnsString(columns[1:]), params)
	pstmt, err := db.DB.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer pstmt.Close()

	result, err := pstmt.Exec(deposit.toTableValues()[1:]...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func UpdateDeposit(deposit *Deposit) error {
	sql := fmt.Sprintf("UPDATE Deposit SET %s WHERE id=?",
		db.ToSettersString(columns[1:]))

	pstmt, err := db.DB.Prepare(sql)
	if err != nil {
		return nil
	}
	defer pstmt.Close()

	values := deposit.toTableValues()[1:]
	values = append(values, deposit.Id)

	_, err = pstmt.Exec(values...)
	return err
}

func DeleteDeposit(id int) error {
	sql := "DELETE FROM Deposit WHERE id=?"

	pstmt, err := db.DB.Prepare(sql)
	if err != nil {
		return nil
	}
	defer pstmt.Close()

	_, err = pstmt.Exec(id)
	return err
}
