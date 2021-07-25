package deposit

import (
	"fmt"

	bk "github.com/zzz404/MoneyGo/internal/bank"
	"github.com/zzz404/MoneyGo/internal/coin"
	mb "github.com/zzz404/MoneyGo/internal/member"
	ut "github.com/zzz404/MoneyGo/internal/utils"
)

type timeDepositController struct {
}

var TimeDepController = new(timeDepositController)

func (c *timeDepositController) Edit(r *ut.HttpRequest, w *ut.HttpResponse) {
	id, isEdit, err := r.GetIntParameter("id", false)
	if w.ResponseForError(err) {
		return
	}
	var td *TimeDeposit
	if isEdit {
		td, err = TimeDepService.Get(id)
		if err == nil && td == nil {
			err = fmt.Errorf("time deposit %d 不存在", id)
		}
		if w.ResponseForError(err) {
			return
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
		coinTypeCode, _, err := r.GetParameter("coinTypeCode", false)
		if w.ResponseForError(err) {
			return
		}
		d := &Deposit{MemberId: memberId, BankId: bankId, CoinTypeCode: coinTypeCode}
		td = &TimeDeposit{Deposit: d}
	}

	data := map[string]interface{}{
		"deposit":      td,
		"members":      mb.Members,
		"banks":        bk.Banks,
		"bankAccounts": bk.BankAccounts,
		"coinTypes":    coin.CoinTypes,
	}
	if isEdit {
		data["id"] = td.Id
	}

	tpl, err := ut.GetTemplate("/timeDepositEdit.html")
	if w.ResponseForError(err) {
		return
	}
	err = tpl.Execute(w, data)

	w.ResponseForError(err)
}

func (c *timeDepositController) Update(r *ut.HttpRequest, w *ut.HttpResponse) {
	deposit, hasId, err := DepController.readDepositFromRequest(r)
	if w.ResponseForError(err) {
		return
	}
	td := &TimeDeposit{Deposit: deposit}

	startDate, _, err := r.GetDateParameter("startDate", true)
	if w.ResponseForError(err) {
		return
	}
	td.StartDate = *startDate

	duration, _, err := r.GetIntParameter("duration", true)
	if w.ResponseForError(err) {
		return
	}
	td.Duration = duration

	td.InterestRate, _, err = r.GetFloatParameter("interestRate", true)
	if w.ResponseForError(err) {
		return
	}

	td.RateTypeCode, _, err = r.GetIntParameter("rateTypeCode", true)
	if w.ResponseForError(err) {
		return
	}

	td.RateTypeCode, _, err = r.GetIntParameter("rateTypeCode", true)
	if w.ResponseForError(err) {
		return
	}

	autoSaveNew, found, err := r.GetBoolParameter("autoSaveNew", false)
	if w.ResponseForError(err) {
		return
	}
	if found {
		td.AutoSaveNew = &autoSaveNew
	}

	if hasId {
		err = TimeDepService.Update(td)
	} else {
		_, err = TimeDepService.Add(td)
	}
	if w.ResponseForError(err) {
		return
	}
	w.Redirect("/static/ReloadOpenerThenClose.html", r)
}

func (c *timeDepositController) List(r *ut.HttpRequest, w *ut.HttpResponse) {
	form := new(QueryForm)
	err := form.ReadParameters(r)
	if w.ResponseForError(err) {
		return
	}

	tds, err := TimeDepService.Query(form)
	if w.ResponseForError(err) {
		return
	}

	totalYearIncome := 0.0
	for _, td := range tds {
		totalYearIncome += td.EspectedYearIncome()
	}

	tpl, err := ut.GetTemplate("/timeDepositList.html")
	if w.ResponseForError(err) {
		return
	}
	err = tpl.Execute(w, map[string]interface{}{
		"form":            form,
		"members":         mb.Members,
		"banks":           bk.Banks,
		"coinTypes":       coin.CoinTypes,
		"tds":             tds,
		"count":           len(tds),
		"totalYearIncome": fmt.Sprintf("%.2f", totalYearIncome),
	})
	w.ResponseForError(err)
}
