package deposit

import (
	"fmt"
	"time"

	bk "github.com/zzz404/MoneyGo/internal/bank"
	"github.com/zzz404/MoneyGo/internal/coin"
	mb "github.com/zzz404/MoneyGo/internal/member"
	"github.com/zzz404/MoneyGo/internal/utils"
)

type DepositType struct {
	Code int
	Name string
}

var DemandDepositType DepositType = DepositType{Code: 1, Name: "活存"}
var TimeDepositType DepositType = DepositType{Code: 2, Name: "定存"}

var DepositTypes = []*DepositType{&DemandDepositType, &TimeDepositType}

func GetDepositTypeByCode(code int) *DepositType {
	switch code {
	case DemandDepositType.Code:
		return &DemandDepositType
	case TimeDepositType.Code:
		return &TimeDepositType
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
	return utils.FormatDate(d.CreatedTime)
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
