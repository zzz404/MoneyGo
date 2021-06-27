package bank

import "github.com/zzz404/MoneyGo/internal/utils"

type bankAccController struct {
}

var BankAccController = new(bankAccController)

func (c *bankAccController) List(r *utils.HttpRequest, w *utils.HttpResponse) {
	tpl, err := utils.GetTemplate("/accountList.html")
	if w.ResponseForError(err) {
		return
	}

	err = tpl.Execute(w, map[string]interface{}{
		"accounts": BankAccounts,
	})
	w.ResponseForError(err)
}

func (c *bankAccController) Edit(r *utils.HttpRequest, w *utils.HttpResponse) {
	tpl, err := utils.GetTemplate("/accountEdit.html")
	if w.ResponseForError(err) {
		return
	}

	data := map[string]interface{}{
		"banks": Banks,
	}
	err = tpl.Execute(w, data)
	w.ResponseForError(err)
}

func (c *bankAccController) Add(r *utils.HttpRequest, w *utils.HttpResponse) {
	account := new(BankAccount)
	var err error

	account.BankId, _, err = r.GetIntParameter("bankId", true)
	if w.ResponseForError(err) {
		return
	}

	account.Account, _, err = r.GetParameter("account", true)
	if w.ResponseForError(err) {
		return
	}

	err = AddBankAccount(account)
	if w.ResponseForError(err) {
		return
	}

	w.Redirect("/static/ReloadOpenerThenClose.html", r)
}

func (c *bankAccController) Delete(r *utils.HttpRequest, w *utils.HttpResponse) {
	account, _, err := r.GetParameter("account", true)
	if w.ResponseJsonError(err) {
		return
	}

	err = DeleteBankAccount(account)
	if w.ResponseJsonError(err) {
		return
	}

	w.WriteJson(true, "", nil)
}
