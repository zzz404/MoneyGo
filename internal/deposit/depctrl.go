package deposit

import (
	"fmt"
	"strconv"

	bk "github.com/zzz404/MoneyGo/internal/bank"
	"github.com/zzz404/MoneyGo/internal/coin"
	mb "github.com/zzz404/MoneyGo/internal/member"
	ut "github.com/zzz404/MoneyGo/internal/utils"
)

func (f *QueryForm) ReadParameters(r *ut.HttpRequest) error {
	memberId, ok, err := r.GetIntParameter("memberId", false)
	if err != nil {
		return err
	}
	if ok {
		f.MemberId = memberId
	}

	f.BankId, _, err = r.GetIntParameter("bankId", false)
	if err != nil {
		return err
	}

	f.TypeCode, _, err = r.GetIntParameter("typeCode", false)
	if err != nil {
		return err
	}

	f.CoinTypeCode, _, err = r.GetParameter("coinTypeCode", false)
	return err
}

type depositController struct {
}

var DepController = new(depositController)

func (c *depositController) List(r *ut.HttpRequest, w *ut.HttpResponse) {
	form := new(QueryForm)
	err := form.ReadParameters(r)
	if w.ResponseForError(err) {
		return
	}

	deposits, err := DepService.Query(form)
	if w.ResponseForError(err) {
		return
	}

	totalTWD := 0.0
	for _, d := range deposits {
		totalTWD += d.TwAmount()
	}

	tpl, err := ut.GetTemplate("/depositList.html")
	if w.ResponseForError(err) {
		return
	}
	err = tpl.Execute(w, map[string]interface{}{
		"form":      form,
		"members":   mb.Members,
		"banks":     bk.Banks,
		"coinTypes": coin.CoinTypes,
		"types":     DepositTypes,
		"deposits":  deposits,
		"count":     len(deposits),
		"totalTWD":  fmt.Sprintf("%.2f", totalTWD),
	})
	w.ResponseForError(err)
}

func (c *depositController) Edit(r *ut.HttpRequest, w *ut.HttpResponse) {
	id, isEdit, err := r.GetIntParameter("id", false)
	if w.ResponseForError(err) {
		return
	}
	var deposit *Deposit
	var td *TimeDeposit
	if isEdit {
		deposit, err = DepService.Get(id)
		if err == nil && deposit == nil {
			err = fmt.Errorf("deposit %d 不存在", id)
		}
		if w.ResponseForError(err) {
			return
		}
		if deposit.TypeCode == TimeDepositType.Code {
			td, err = TimeDepService.GetTd(deposit)
			if w.ResponseForError(err) {
				return
			}
		}
	} else {
		memberId, _, err := r.GetIntParameter("memberId", false)
		if w.ResponseForError(err) {
			return
		}
		bankId, _, err := r.GetIntParameter("bankId", false)
		if w.ResponseForError(err) {
			return
		}
		typeCode, _, err := r.GetIntParameter("typeCode", false)
		if w.ResponseForError(err) {
			return
		}
		coinTypeCode, _, err := r.GetParameter("coinTypeCode", false)
		if w.ResponseForError(err) {
			return
		}
		deposit = &Deposit{MemberId: memberId, BankId: bankId, TypeCode: typeCode, CoinTypeCode: coinTypeCode}
	}
	if td == nil {
		td = &TimeDeposit{Deposit: deposit}
	}

	radioFixed := ut.Radio{
		Text:    InterestRateType_FIXED.Name,
		Value:   strconv.Itoa(InterestRateType_FIXED.Code),
		Checked: (td.RateTypeCode != nil && *td.RateTypeCode == InterestRateType_FIXED.Code),
	}
	radioVariable := ut.Radio{
		Text:    InterestRateType_VARIABLE.Name,
		Value:   strconv.Itoa(InterestRateType_VARIABLE.Code),
		Checked: (td.RateTypeCode != nil && *td.RateTypeCode == InterestRateType_VARIABLE.Code),
	}

	radioYes := ut.Radio{
		Text:    "是",
		Value:   "true",
		Checked: (td.AutoSaveNew != nil && *td.AutoSaveNew),
	}
	radioNo := ut.Radio{
		Text:    "否",
		Value:   "false",
		Checked: (td.AutoSaveNew != nil && !*td.AutoSaveNew),
	}

	data := map[string]interface{}{
		"deposit":           td,
		"members":           mb.Members,
		"banks":             bk.Banks,
		"bankAccounts":      bk.BankAccounts,
		"depositTypes":      DepositTypes,
		"coinTypes":         coin.CoinTypes,
		"interestRateTypes": []*ut.Radio{&radioFixed, &radioVariable},
		"timeDepCode":       TimeDepositType.Code,
		"autoSaveNews":      []*ut.Radio{&radioYes, &radioNo},
	}
	if isEdit {
		data["id"] = deposit.Id
	}

	tpl, err := ut.GetTemplate("/depositEdit.html")
	if w.ResponseForError(err) {
		return
	}
	err = tpl.Execute(w, data)

	w.ResponseForError(err)
}

func (c *depositController) readDepositFromRequest(r *ut.HttpRequest) (deposit *Deposit, hasId bool, err error) {
	deposit = &Deposit{}

	id, hasId, err := r.GetIntParameter("id", false)
	if err != nil {
		return
	}
	if hasId {
		deposit.Id = id
	}

	deposit.MemberId, _, err = r.GetIntParameter("memberId", true)
	if err != nil {
		return
	}

	deposit.BankId, _, err = r.GetIntParameter("bankId", true)
	if err != nil {
		return
	}

	deposit.BankAccount, _, err = r.GetParameter("bankAccount", false)
	if err != nil {
		return
	}

	deposit.TypeCode, _, err = r.GetIntParameter("typeCode", true)
	if err != nil {
		return
	}

	deposit.Amount, _, err = r.GetFloatParameter("amount", true)
	if err != nil {
		return
	}

	deposit.CoinTypeCode, _, err = r.GetParameter("coinTypeCode", true)
	if err != nil {
		return
	}

	return
}

func (c *depositController) readTimeDepositFromRequest(td *TimeDeposit, r *ut.HttpRequest) error {
	startDate, err := r.GetDatePointerParameter("startDate", false)
	if err != nil {
		return err
	}
	td.StartDate = startDate

	endDate, err := r.GetDatePointerParameter("endDate", false)
	if err != nil {
		return err
	}
	td.EndDate = endDate

	duration, err := r.GetIntPointerParameter("duration", false)
	if err != nil {
		return err
	}
	td.Duration = duration

	interestRate, err := r.GetFloatPointerParameter("interestRate", false)
	if err != nil {
		return err
	}
	td.InterestRate = interestRate

	rateTypeCode, err := r.GetIntPointerParameter("rateTypeCode", false)
	if err != nil {
		return err
	}
	td.RateTypeCode = rateTypeCode

	autoSaveNew, found, err := r.GetBoolParameter("autoSaveNew", false)
	if err != nil {
		return err
	}
	if found {
		td.AutoSaveNew = &autoSaveNew
	}

	return nil
}

func (c *depositController) Update(r *ut.HttpRequest, w *ut.HttpResponse) {
	deposit, hasId, err := c.readDepositFromRequest(r)
	if w.ResponseForError(err) {
		return
	}

	var td *TimeDeposit
	if deposit.TypeCode == TimeDepositType.Code {
		td = &TimeDeposit{Deposit: deposit}
		err := c.readTimeDepositFromRequest(td, r)
		if w.ResponseForError(err) {
			return
		}
	}

	if hasId {
		if td != nil {
			err = TimeDepService.Update(td)
		} else {
			err = DepService.Update(deposit)
		}
		if w.ResponseForError(err) {
			return
		}
	} else {
		if td != nil {
			_, err = TimeDepService.Add(td)
		} else {
			_, err = DepService.Add(deposit)
		}
		if w.ResponseForError(err) {
			return
		}
	}
	w.Redirect("/static/ReloadOpenerThenClose.html", r)
}

func (c *depositController) Delete(r *ut.HttpRequest, w *ut.HttpResponse) {
	id, _, err := r.GetIntParameter("id", true)
	if w.ResponseJsonError(err) {
		return
	}

	err = DepService.Delete(id)
	if w.ResponseJsonError(err) {
		return
	}
	w.WriteJson(true, "", nil)
}
