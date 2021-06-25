package deposit

import (
	"fmt"
	"time"

	bk "github.com/zzz404/MoneyGo/internal/bank"
	"github.com/zzz404/MoneyGo/internal/coin"
	"github.com/zzz404/MoneyGo/internal/db"
	mb "github.com/zzz404/MoneyGo/internal/member"
)

type DepositType struct {
	Code int
	Name string
}

var DemandDeposit DepositType = DepositType{Code: 1, Name: "活存"}
var TimeDeposit DepositType = DepositType{Code: 2, Name: "定存"}

var DepositTypes = []*DepositType{&DemandDeposit, &TimeDeposit}

func GetDepositTypeByCode(code int) *DepositType {
	switch code {
	case DemandDeposit.Code:
		return &DemandDeposit
	case TimeDeposit.Code:
		return &TimeDeposit
	}
	panic(fmt.Errorf("DepositType %d 不存在", code))
}

type Deposit struct {
	Id                int
	MemberId          int
	BankId            int
	BankAccount       string
	TypeCode          int
	Amount            float64
	CoinTypeCode      string
	CreatedTime       time.Time
	ExRateWhenCreated float64

	_Type     *DepositType
	_Member   *mb.Member
	_Bank     *bk.Bank
	_CoinType *coin.CoinType
}

func (d *Deposit) Type() *DepositType {
	if d._Type == nil {
		d._Type = GetDepositTypeByCode(d.TypeCode)
	}
	return d._Type
}

func (d *Deposit) Member() *mb.Member {
	if d._Member == nil {
		d._Member = mb.GetMember(d.MemberId)
	}
	return d._Member
}

func (d *Deposit) Bank() *bk.Bank {
	if d._Bank == nil {
		d._Bank = bk.GetBank(d.BankId)
	}
	return d._Bank
}

func (d *Deposit) CoinType() *coin.CoinType {
	if d._CoinType == nil {
		d._CoinType = coin.GetCoinTypeByCode(d.CoinTypeCode)
	}
	return d._CoinType
}

func (d *Deposit) CreatedTimeString() string {
	return d.CreatedTime.Format("2006-01-02 15:04:05")
}

func (d *Deposit) AmountString() string {
	return fmt.Sprintf("%.2f", d.Amount)
}

func (d *Deposit) TwAmount() float64 {
	return d.Amount * d.CoinType().ExRate
}

func (d *Deposit) TwAmountString() string {
	return fmt.Sprintf("%.2f", d.TwAmount())
}

var columnsForUpdate = []string{"bankId", "bankAccount", "type", "amount", "coinType"}
var columnsForInsert = append(columnsForUpdate, "memberId", "exRateWhenCreated")
var columnsForQuery = append(columnsForInsert, "id", "createdTime")

func (d *Deposit) toValuesOfUpdate() []interface{} {
	return []interface{}{d.BankId, d.BankAccount, d.TypeCode, d.Amount, d.CoinTypeCode}
}

func (d *Deposit) toValuesOfInsert() []interface{} {
	return append(d.toValuesOfUpdate(), d.MemberId, d.ExRateWhenCreated)
}

func (d *Deposit) toValuesOfQuery() []interface{} {
	return []interface{}{
		&d.BankId, &d.TypeCode, &d.Amount, &d.CoinTypeCode,
		&d.MemberId, &d.ExRateWhenCreated,
		&d.Id, &d.CreatedTime}
}

type QueryForm struct {
	MemberId     int
	BankId       int
	TypeCode     int
	CoinTypeCode string
}

func QueryDeposits(form *QueryForm) ([]*Deposit, error) {
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

		err := rows.Scan(deposit.toValuesOfQuery()...)
		if err != nil {
			return nil, err
		}

		deposits = append(deposits, deposit)
	}
	return deposits, nil
}

func GetDeposit(id int) (*Deposit, error) {
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
		err := rows.Scan(deposit.toValuesOfQuery()...)
		if err != nil {
			return nil, err
		}
	}
	return deposit, nil
}

func AddDeposit(deposit *Deposit) (int, error) {
	params, err := db.ToSqlParams(len(columnsForInsert))
	if err != nil {
		return 0, err
	}
	sql := fmt.Sprintf("INSERT INTO Deposit (%s) VALUES (%s)",
		db.ToColumnsString(columnsForInsert), params)

	deposit.ExRateWhenCreated = deposit.CoinType().ExRate
	result, err := db.ExecuteSql(sql, deposit.toValuesOfInsert()...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func UpdateDeposit(deposit *Deposit) error {
	sql := fmt.Sprintf("UPDATE Deposit SET %s WHERE id=?",
		db.ToSettersString(columnsForUpdate))
	values := append(deposit.toValuesOfUpdate(), deposit.Id)

	_, err := db.ExecuteSql(sql, values...)
	return err
}

func DeleteDeposit(id int) error {
	sql := "DELETE FROM Deposit WHERE id=?"
	_, err := db.ExecuteSql(sql, id)
	return err
}

func QueryTotalTWD() (map[int]float64, float64, error) {
	m := map[int]float64{}
	sql := `SELECT d.memberId, sum(d.amount*c.exchangeRate) AS totalTWD
            FROM Deposit d, CoinType c 
            WHERE d.coinType=c.code
            GROUP BY d.memberId`
	rows, err := db.DB.Query(sql)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

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
