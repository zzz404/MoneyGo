package deposit

import (
	"fmt"

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

	tpl, err := ut.GetTemplate("/depositList.html")
	if w.ResponseForError(err) {
		return
	}
	deposits, err := DepService.Query(form)
	if w.ResponseForError(err) {
		return
	}

	totalTWD := 0.0
	for _, d := range deposits {
		totalTWD += d.Amount * d.CoinType().ExRate
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
	tpl, err := ut.GetTemplate("/depositEdit.html")
	if w.ResponseForError(err) {
		return
	}
	id, isEdit, err := r.GetIntParameter("id", false)
	if w.ResponseForError(err) {
		return
	}
	var deposit *Deposit
	if isEdit {
		deposit, err = DepService.Get(id)
		if err == nil && deposit == nil {
			err = fmt.Errorf("deposit %d 不存在", id)
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

	data := map[string]interface{}{
		"deposit":      deposit,
		"members":      mb.Members,
		"banks":        bk.Banks,
		"bankAccounts": bk.BankAccounts,
		"depositTypes": DepositTypes,
		"coinTypes":    coin.CoinTypes,
	}
	if isEdit {
		data["id"] = deposit.Id
	}
	err = tpl.Execute(w, data)

	w.ResponseForError(err)
}

func (c *depositController) Update(r *ut.HttpRequest, w *ut.HttpResponse) {
	deposit := &Deposit{}

	id, hasId, err := r.GetIntParameter("id", false)
	if w.ResponseForError(err) {
		return
	}
	if hasId {
		deposit.Id = id
	}

	deposit.MemberId, _, err = r.GetIntParameter("memberId", true)
	if w.ResponseForError(err) {
		return
	}

	deposit.BankId, _, err = r.GetIntParameter("bankId", true)
	if w.ResponseForError(err) {
		return
	}

	deposit.BankAccount, _, err = r.GetParameter("bankAccount", true)
	if w.ResponseForError(err) {
		return
	}

	deposit.TypeCode, _, err = r.GetIntParameter("typeCode", true)
	if w.ResponseForError(err) {
		return
	}

	deposit.Amount, _, err = r.GetFloatParameter("amount", true)
	if w.ResponseForError(err) {
		return
	}

	deposit.CoinTypeCode, _, err = r.GetParameter("coinTypeCode", true)
	if w.ResponseForError(err) {
		return
	}

	if hasId {
		err = DepService.Update(deposit)
		if w.ResponseForError(err) {
			return
		}
	} else {
		_, err = DepService.Add(deposit)
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
