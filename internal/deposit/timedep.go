package deposit

import (
	"fmt"
	"strconv"
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

type AutoSaveNewType struct {
	Code int
	Name string
}

type TimeDeposit struct {
	*Deposit
	StartDate    *time.Time
	EndDate      *time.Time
	Duration     *int     // 存幾個月
	InterestRate *float64 // 年利率
	RateTypeCode *int     // 固定或機動
	AutoSaveNew  *bool    // 自動轉存
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

func (td *TimeDeposit) InterestRatePercentString() string {
	if td.InterestRate == nil {
		return ""
	} else {
		percent := *td.InterestRate
		return fmt.Sprintf("%.2f", percent)
	}
}

func (td *TimeDeposit) RateTypeString() string {
	if td.RateTypeCode == nil {
		return ""
	} else {
		return GetInterestRateTypeByCode(*td.RateTypeCode).Name
	}
}

func (td *TimeDeposit) AutoSaveNewString() string {
	if td.AutoSaveNew == nil {
		return ""
	} else if *td.AutoSaveNew {
		return "是"
	} else {
		return "否"
	}
}

func (td *TimeDeposit) EspectedYearIncome() *float64 {
	if td.InterestRate == nil {
		return nil
	} else {
		income := td.TwAmount() * (*td.InterestRate)
		return &income
	}
}

func (td *TimeDeposit) EspectedYearIncomeString() string {
	income := td.EspectedYearIncome()
	if income == nil {
		return ""
	} else {
		return fmt.Sprintf("%g", income) + "%"
	}
}

func (td *TimeDeposit) DurationString() string {
	duration := td.Duration
	if duration == nil {
		return ""
	} else {
		return strconv.Itoa(*duration)
	}
}
