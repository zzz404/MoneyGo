package deposit

import (
	"fmt"

	"github.com/zzz404/MoneyGo/internal/db"
)

type timeDepositService struct {
	columns []string
}

var Service = &timeDepositService{
	columns: []string{"depId", "startDate", "endDate", "interestRate", "rateTypeCode", "autoSaveNew"},
}

func (s *timeDepositService) toValues(td *TimeDeposit) []interface{} {
	return []interface{}{td.deposit.Id, td.StartDate, td.EndDate, td.InterestRate, td.RateTypeCode, td.AutoSaveNew}
}

func (s *timeDepositService) Add(td *TimeDeposit) (int, error) {
	exe := &db.SqlExecuter{}
	defer exe.Close()

	err := exe.StartTx()
	if err != nil {
		return 0, err
	}

	id, err := addDeposit(td.deposit, exe)
	if err != nil {
		return 0, err
	}

	err = addTimeDeposit(td, exe)
	if err != nil {
		return 0, err
	}

	err = exe.Commit()
	if err != nil {
		return 0, err
	}

	return id, err
}

// var t_columns = []string{"depId", "startDate", "endDate", "interestRate", "rateTypeCode", "autoSaveNew"}

// func (d *TimeDeposit) t_toValues() []interface{} {
// 	return []interface{}{d.deposit.Id, d.StartDate, d.EndDate, d.InterestRate, d.RateTypeCode, d.AutoSaveNew}
// }

func AddTimeDeposit(td *TimeDeposit) (int, error) {
	exe := &db.SqlExecuter{}
	defer exe.Close()

	err := exe.StartTx()
	if err != nil {
		return 0, err
	}

	id, err := addDeposit(td.deposit, exe)
	if err != nil {
		return 0, err
	}

	err = addTimeDeposit(td, exe)
	if err != nil {
		return 0, err
	}

	err = exe.Commit()
	if err != nil {
		return 0, err
	}

	return id, err
}

func addTimeDeposit(td *TimeDeposit, exe *db.SqlExecuter) error {
	params, err := db.ToSqlParams(len(t_columns))
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("INSERT INTO TimeDeposit (%s) VALUES (%s)",
		db.ToColumnsString(t_columns), params)

	_, err = exe.ExecuteSql(sql, td.t_toValues()...)
	return err
}

func UpdateTimeDeposit(deposit *TimeDeposit) error {
	exe := &db.SqlExecuter{}
	defer exe.Close()

	err := exe.StartTx()
	if err != nil {
		return err
	}

	err = updateDeposit(deposit.deposit, exe)
	if err != nil {
		return err
	}

	err = updateTimeDeposit(deposit, exe)
	if err != nil {
		return err
	}

	return exe.Commit()
}

func updateTimeDeposit(deposit *TimeDeposit, exe *db.SqlExecuter) error {
	sql := fmt.Sprintf("UPDATE TimeDeposit SET %s WHERE depId=?",
		db.ToSettersString(t_columns[1:]))
	values := deposit.t_toValues()
	values = append(values[1:], values[0])

	_, err := exe.ExecuteSql(sql, values...)
	return err
}

func GetTimeDeposit(id int) (*Deposit, error) {
	sql := "SELECT " + db.ToColumnsString(columnsForQuery) + " FROM Deposit WHERE id=?"
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
		err := deposit.loadFromRows(rows)
		if err != nil {
			return nil, err
		}
	}
	return deposit, nil
}

func QueryTimeDeposits(form *QueryForm) ([]*Deposit, error) {
	sb := &db.SqlBuilder{}
	sb.Columns = columnsForQuery
	sb.AddTable("Deposit").SetColumns(columnsForQuery)
	if form.MemberId > 0 {
		sb.AddCondition("memberId=?", form.MemberId)
	}
	if form.BankId > 0 {
		sb.AddCondition("bankId=?", form.BankId)
	}
	if form.TypeCode > 0 {
		sb.AddCondition("type=?", form.TypeCode)
	}
	if form.CoinTypeCode != "" {
		sb.AddCondition("coinType=?", form.CoinTypeCode)
	}
	sb.AddOrderBy("bankId ASC").AddOrderBy("id DESC")
	sql := sb.BuildSql()

	rows, err := db.DB.Query(sql, sb.Variables...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deposits []*Deposit
	for rows.Next() {
		deposit := &Deposit{}
		err := deposit.loadFromRows(rows)
		if err != nil {
			return nil, err
		}

		deposits = append(deposits, deposit)
	}
	return deposits, nil
}
