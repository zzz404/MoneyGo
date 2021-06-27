package deposit

import (
	"database/sql"
	"fmt"

	"github.com/zzz404/MoneyGo/internal/db"
	"github.com/zzz404/MoneyGo/internal/utils"
)

type depositService struct {
	columnsForUpdate []string
	columnsForInsert []string
	columnsForQuery  []string
}

var DepService = func() *depositService {
	var columnsForUpdate = []string{"bankId", "bankAccount", "type", "amount", "coinType"}
	var columnsForInsert = append(columnsForUpdate, "memberId", "exRateWhenCreated")
	var columnsForQuery = append(columnsForInsert, "id", "createdTime")
	return &depositService{
		columnsForUpdate: columnsForUpdate,
		columnsForInsert: columnsForInsert,
		columnsForQuery:  columnsForQuery,
	}
}()

func (s *depositService) toValuesOfUpdate(d *Deposit) []interface{} {
	return []interface{}{d.BankId, d.BankAccount, d.TypeCode, d.Amount, d.CoinTypeCode}
}

func (s *depositService) toValuesOfInsert(d *Deposit) []interface{} {
	return append(s.toValuesOfUpdate(d), d.MemberId, d.CoinType().ExRate)
}

func (s *depositService) loadFromRows(d *Deposit, rows *sql.Rows) error {
	var bankAccount sql.NullString
	err := rows.Scan(&d.BankId, &bankAccount, &d.TypeCode, &d.Amount, &d.CoinTypeCode,
		&d.MemberId, &d.ExRateWhenCreated,
		&d.Id, &d.CreatedTime)
	if err != nil {
		return err
	}
	if bankAccount.Valid {
		d.BankAccount = bankAccount.String
	}
	return nil
}

type QueryForm struct {
	MemberId     int
	BankId       int
	TypeCode     int
	CoinTypeCode string
}

func (f *QueryForm) SetToSqlBuilder(sb *db.SqlBuilder, prefix string) {
	if prefix != "" {
		prefix = prefix + "."
	}
	if f.MemberId > 0 {
		sb.AddCondition(prefix+"memberId=?", f.MemberId)
	}
	if f.BankId > 0 {
		sb.AddCondition(prefix+"bankId=?", f.BankId)
	}
	if f.TypeCode > 0 {
		sb.AddCondition(prefix+"type=?", f.TypeCode)
	}
	if f.CoinTypeCode != "" {
		sb.AddCondition(prefix+"coinType=?", f.CoinTypeCode)
	}
}

func (s *depositService) Query(form *QueryForm) (ds []*Deposit, err error) {
	sb := &db.SqlBuilder{}
	sb.AddTable("Deposit").SetColumns(s.columnsForQuery)
	form.SetToSqlBuilder(sb, "")
	sb.AddOrderBy("bankId ASC").AddOrderBy("id DESC")
	sql := sb.BuildSql()

	rows, err := db.DB.Query(sql, sb.Variables...)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = utils.CombineError(err, rows.Close())
	}()

	for rows.Next() {
		deposit := &Deposit{}
		err = s.loadFromRows(deposit, rows)
		if err != nil {
			return
		}
		ds = append(ds, deposit)
	}
	return
}

func (s *depositService) Get(id int) (d *Deposit, err error) {
	sql := "SELECT " + db.ToColumnsString(s.columnsForQuery) + " FROM Deposit WHERE id=?"
	rows, err := db.DB.Query(sql, id)
	if err != nil {
		return
	}
	defer func() {
		err = utils.CombineError(err, rows.Close())
	}()

	for rows.Next() {
		if d == nil {
			d = &Deposit{}
		} else {
			return nil, fmt.Errorf("Deposit id %d 不只一個!?", id)
		}
		err = s.loadFromRows(d, rows)
		if err != nil {
			return
		}
	}
	return
}

func (s *depositService) Add(deposit *Deposit) (int, error) {
	return s.add(deposit, db.DB)
}

func (s *depositService) add(dep *Deposit, exe db.SqlExecuter) (int, error) {
	params, err := db.ToSqlParams(len(s.columnsForInsert))
	if err != nil {
		return 0, err
	}
	sql := fmt.Sprintf("INSERT INTO Deposit (%s) VALUES (%s)",
		db.ToColumnsString(s.columnsForInsert), params)

	result, err := exe.Exec(sql, s.toValuesOfInsert(dep)...)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (s *depositService) Update(deposit *Deposit) error {
	return s.update(deposit, db.DB)
}

func (s *depositService) update(dep *Deposit, exe db.SqlExecuter) error {
	sql := fmt.Sprintf("UPDATE Deposit SET %s WHERE id=?",
		db.ToSettersString(s.columnsForUpdate))
	values := append(s.toValuesOfUpdate(dep), dep.Id)

	_, err := exe.Exec(sql, values...)
	return err
}

func (s *depositService) Delete(id int) error {
	sql := "DELETE FROM Deposit WHERE id=?"
	_, err := db.DB.Exec(sql, id)
	return err
}

func (s *depositService) QueryTotalTWD() (map[int]float64, float64, error) {
	m := map[int]float64{}
	sql := `SELECT d.memberId, sum(d.amount*c.exchangeRate) AS totalTWD
            FROM Deposit d, CoinType c 
            WHERE d.coinType=c.code
            GROUP BY d.memberId`
	rows, err := db.DB.Query(sql)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		err = utils.CombineError(err, rows.Close())
	}()

	var memberId int
	var memberTotal float64
	var total float64 = 0

	for rows.Next() {
		err := rows.Scan(&memberId, &memberTotal)
		if err != nil {
			return nil, 0, err
		}
		m[memberId] = memberTotal
		total += memberTotal
	}
	return m, total, nil
}
