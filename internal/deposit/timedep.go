package deposit

import (
	"fmt"
	"time"

	"github.com/zzz404/MoneyGo/internal/utils"
)

type InterestRateType struct {
	Code int
	Name string
}

var InterestRateType_FIXED InterestRateType = InterestRateType{Code: 1, Name: "固定"}
var InterestRateType_VARIABLE InterestRateType = InterestRateType{Code: 2, Name: "機動"}

var InterestRateTypes = []*InterestRateType{&InterestRateType_FIXED, &InterestRateType_VARIABLE}

func GetInterestRateTypeByCode(code int) *InterestRateType {
	switch code {
	case InterestRateType_FIXED.Code:
		return &InterestRateType_FIXED
	case InterestRateType_VARIABLE.Code:
		return &InterestRateType_VARIABLE
	}
	panic(fmt.Errorf("InterestRateType %d 不存在", code))
}

type TimeDeposit struct {
	*Deposit
	StartDate    time.Time
	EndDate      time.Time
	InterestRate float64
	RateTypeCode int
	AutoSaveNew  int
}

func NewTimeDeposit() *TimeDeposit {
	return &TimeDeposit{Deposit: new(Deposit)}
}

func (td *TimeDeposit) StartDateString() string {
	return utils.FormatDate(td.StartDate)
}

func (td *TimeDeposit) EndDateString() string {
	return utils.FormatDate(td.EndDate)
}

func (td *TimeDeposit) RateType() *InterestRateType {
	return GetInterestRateTypeByCode(td.RateTypeCode)
}

func (td *TimeDeposit) AutoSaveNewString() string {
	if td.AutoSaveNew == 0 {
		return "否"
	} else {
		return "是"
	}
}

func (td *TimeDeposit) EspectedYearIncome() float64 {
	return td.TwAmount() * td.InterestRate
}
