package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zzz404/MoneyGo/internal/db"
	dp "github.com/zzz404/MoneyGo/internal/deposit"
)

func memberList(r *HttpRequest, w *HttpResponse) {
	tpl, err := getTemplate("/memberList.html")
	if w.responseForError(err) {
		return
	}
	members, err := db.QueryMembers()
	if w.responseForError(err) {
		return
	}
	err = tpl.Execute(w, map[string]interface{}{
		"members": members,
	})
	w.responseForError(err)
}

func depositList(r *HttpRequest, w *HttpResponse) {
	memberId, _, err := r.getIntParameter("memberId", true)
	if w.responseForError(err) {
		return
	}
	tpl, err := getTemplate("/depositList.html")
	if w.responseForError(err) {
		return
	}
	deposits, err := dp.QueryDeposits(memberId)
	if w.responseForError(err) {
		return
	}
	err = tpl.Execute(w, map[string]interface{}{
		"memberId": memberId,
		"deposits": deposits,
	})
	w.responseForError(err)
}

func depositEdit(r *HttpRequest, w *HttpResponse) {
	tpl, err := getTemplate("/depositEdit.html")
	if w.responseForError(err) {
		return
	}
	id, found, err := r.getIntParameter("id", false)
	if w.responseForError(err) {
		return
	}
	var deposit *dp.Deposit
	if found {
		deposit, err = dp.GetDeposit(id)
		if err == nil && deposit == nil {
			err = fmt.Errorf("deposit %d 不存在", id)
		}
		if w.responseForError(err) {
			return
		}
	} else {
		memberId, _, err := r.getIntParameter("memberId", true)
		if w.responseForError(err) {
			return
		}
		deposit = &dp.Deposit{MemberId: memberId}
	}
	err = tpl.Execute(w, deposit)
	w.responseForError(err)
}

func depositUpdate(r *HttpRequest, w *HttpResponse) {
	deposit := &dp.Deposit{}

	id, hasId, err := r.getIntParameter("id", false)
	if w.responseForError(err) {
		return
	}
	if hasId {
		deposit.Id = id
	}

	deposit.MemberId, _, err = r.getIntParameter("memberId", true)
	if w.responseForError(err) {
		return
	}

	deposit.BankId, _, err = r.getIntParameter("bankId", true)
	if w.responseForError(err) {
		return
	}

	typeCode, _, err := r.getIntParameter("typeCode", true)
	if w.responseForError(err) {
		return
	}
	deposit.Type, err = dp.GetDepositTypeByCode(typeCode)
	if w.responseForError(err) {
		return
	}

	deposit.Amount, _, err = r.getFloatParameter("amount", true)
	if w.responseForError(err) {
		return
	}

	coinTypeCode := r.URL.Query().Get("coinTypeCode")
	deposit.CoinType, err = dp.GetCoinTypeByCode(coinTypeCode)
	if w.responseForError(err) {
		return
	}

	if hasId {
		err = dp.UpdateDeposit(deposit)
		if w.responseForError(err) {
			return
		}
	} else {
		id, err = dp.AddDeposit(deposit)
		if w.responseForError(err) {
			return
		}
	}

	url := fmt.Sprintf("/depositEdit?id=%d", id)
	http.Redirect(w, r.Request, url, http.StatusSeeOther)
}

func Start() {
	http.HandleFunc("/", makeHandler(memberList))
	http.HandleFunc("/depositList", makeHandler(depositList))
	http.HandleFunc("/depositEdit", makeHandler(depositEdit))
	http.HandleFunc("/depositUpdate", makeHandler(depositUpdate))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
