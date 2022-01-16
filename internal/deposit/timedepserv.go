package deposit

import (
	"database/sql"
	"fmt"

	"github.com/zzz404/MoneyGo/internal/db"
	"github.com/zzz404/MoneyGo/internal/utils"
)

type timeDepositService struct {
	columns []string
}

var TimeDepService = &timeDepositService{
	columns: []string{"startDate", "endDate", "duration", "interestRate", "rateTypeCode", "autoSaveNew"},
}

func (s *timeDepositService) toValues(d *TimeDeposit) []interface{} {
	return []interface{}{d.StartDate, d.EndDate, d.Duration, d.InterestRate, d.RateTypeCode, d.AutoSaveNew}
}

func (s *timeDepositService) loadFromRows(td *TimeDeposit, rows *sql.Rows) error {
	//var autoSaveNew sql.NullBool
	// err := rows.Scan(&td.StartDate, &td.EndDate, &td.InterestRate, &td.RateTypeCode, &autoSaveNew)
	// if err != nil {
	// 	return err
	// }
	// if autoSaveNew.Valid {
	// 	td.AutoSaveNew = &autoSaveNew.Bool
	// }
	// return nil
	return rows.Scan(&td.StartDate, &td.EndDate, &td.Duration, &td.InterestRate, &td.RateTypeCode, &td.AutoSaveNew)
}

func (s *timeDepositService) Add(td *TimeDeposit) (id int, err error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = db.CommitOrRollback(tx, err)
	}()

	id, err = DepService.add(td.Deposit, tx)
	if err != nil {
		return
	}

	err = s.add(td, tx)
	return
}

func (s *timeDepositService) add(td *TimeDeposit, exe db.SqlExecuter) error {
	params, err := db.ToSqlParams(len(s.columns) + 1)
	if err != nil {
		return err
	}
	columns := append(s.columns, "depId")
	sql := fmt.Sprintf("INSERT INTO TimeDeposit (%s) VALUES (%s)",
		db.ToColumnsString(columns), params)

	values := append(s.toValues(td), td.Id)
	_, err = exe.Exec(sql, values...)
	return err
}

func (s *timeDepositService) Update(deposit *TimeDeposit) (err error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = db.CommitOrRollback(tx, err)
	}()

	err = DepService.update(deposit.Deposit, tx)
	if err != nil {
		return
	}

	rowsAffected, err := s.update(deposit, tx)
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = s.add(deposit, tx)
	}

	return
}

func (s *timeDepositService) update(td *TimeDeposit, exe db.SqlExecuter) (int, error) {
	sql := fmt.Sprintf("UPDATE TimeDeposit SET %s WHERE depId=?",
		db.ToSettersString(s.columns))
	values := append(s.toValues(td), td.Id)

	result, err := exe.Exec(sql, values...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (s *timeDepositService) GetTd(dep *Deposit) (td *TimeDeposit, err error) {
	sql := "SELECT " + db.ToColumnsString(s.columns) + " FROM TimeDeposit WHERE depId=?"
	rows, err := db.DB.Query(sql, dep.Id)
	if err != nil {
		return
	}
	defer func() {
		err = utils.CombineError(err, rows.Close())
	}()

	for rows.Next() {
		if td == nil {
			td = &TimeDeposit{Deposit: dep}
			err = s.loadFromRows(td, rows)
			if err != nil {
				return
			}
		} else {
			return nil, fmt.Errorf("TimeDeposit id %d 不只一個!?", dep.Id)
		}
	}
	return
}

func (s *timeDepositService) Query(form *QueryForm) (tds []*TimeDeposit, err error) {
	sb := &db.SqlBuilder{}

	sb.AddTable("Deposit d").AddTable("TimeDeposit td").SetColumns(
		[]string{"d.memberId", "d.bankId", "d.amount", "d.coinType",
			"td.startDate", "td.endDate", "interestRate", "rateTypeCode", "autoSaveNew"},
	)
	var loadRows = func(td *TimeDeposit, rows *sql.Rows) error {
		// var autoSaveNew sql.NullBool
		// err := rows.Scan(&td.MemberId, &td.BankId, &td.Amount, &td.CoinTypeCode,
		// 	&td.StartDate, &td.EndDate, &td.InterestRate, &td.RateTypeCode, &autoSaveNew)
		// if err != nil {
		// 	return err
		// }
		// if autoSaveNew.Valid {
		// 	td.AutoSaveNew = &autoSaveNew.Bool
		// }
		// return nil
		return rows.Scan(&td.MemberId, &td.BankId, &td.Amount, &td.CoinTypeCode,
			td.StartDate, td.Duration, td.InterestRate, td.RateTypeCode, td.AutoSaveNew)
	}

	sb.AddCondition("td.depId=d.id", nil)
	form.SetToSqlBuilder(sb, "d")
	sb.AddOrderBy("td.endDate DESC")
	sql := sb.BuildSql()

	rows, err := db.DB.Query(sql, sb.Variables...)
	if err != nil {
		return
	}
	defer func() {
		err = utils.CombineError(err, rows.Close())
	}()

	for rows.Next() {
		td := NewTimeDeposit()
		err = loadRows(td, rows)
		if err != nil {
			return
		}
		tds = append(tds, td)
	}
	return
}
