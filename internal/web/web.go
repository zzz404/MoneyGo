package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zzz404/MoneyGo/internal/db"
	dp "github.com/zzz404/MoneyGo/internal/deposit"
)

func memberList(w http.ResponseWriter, r *http.Request) {
	tpl, err := getTemplate("/memberList.html")
	if responseForError(err, w) {
		return
	}
	members, err := db.QueryMembers()
	if responseForError(err, w) {
		return
	}
	err = tpl.Execute(w, map[string]interface{}{
		"members": members,
	})
	responseForError(err, w)
}

func depositList(w http.ResponseWriter, r *http.Request) {
	memberId, _, err := getIntParameter(r, "memberId", true)
	if responseForError(err, w) {
		return
	}
	tpl, err := getTemplate("/depositList.html")
	if responseForError(err, w) {
		return
	}
	deposits, err := dp.QueryDeposits(memberId)
	if responseForError(err, w) {
		return
	}
	err = tpl.Execute(w, map[string]interface{}{
		"memberId": memberId,
		"deposits": deposits,
	})
	responseForError(err, w)
}

func depositEdit(w http.ResponseWriter, r *http.Request) {
	tpl, err := getTemplate("/depositEdit.html")
	if responseForError(err, w) {
		return
	}
	id, found, err := getIntParameter(r, "id", false)
	if responseForError(err, w) {
		return
	}
	var deposit *dp.Deposit
	if found {
		deposit, err = dp.GetDeposit(id)
		if err == nil && deposit == nil {
			err = fmt.Errorf("deposit %d 不存在", id)
		}
		if responseForError(err, w) {
			return
		}
	} else {
		memberId, _, err := getIntParameter(r, "memberId", true)
		if responseForError(err, w) {
			return
		}
		deposit = &dp.Deposit{MemberId: memberId}
	}
	err = tpl.Execute(w, deposit)
	responseForError(err, w)
}

func depositUpdate(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
	}
	deposit := &dp.Deposit{}

	id, hasId, err := getIntParameter(r, "id", false)
	if responseForError(err, w) {
		return
	}
	if hasId {
		deposit.Id = id
	}

	deposit.MemberId, _, err = getIntParameter(r, "memberId", true)
	if responseForError(err, w) {
		return
	}

	deposit.BankId, _, err = getIntParameter(r, "bankId", true)
	if responseForError(err, w) {
		return
	}

	typeCode, _, err := getIntParameter(r, "typeCode", true)
	if responseForError(err, w) {
		return
	}
	deposit.Type, err = dp.GetDepositTypeByCode(typeCode)
	if responseForError(err, w) {
		return
	}

	deposit.Amount, _, err = getFloatParameter(r, "amount", true)
	if responseForError(err, w) {
		return
	}

	coinTypeCode := r.URL.Query().Get("coinTypeCode")
	deposit.CoinType, err = dp.GetCoinTypeByCode(coinTypeCode)
	responseForError(err, w)

    r.Response.
    w.
}

func Start() {
	http.HandleFunc("/", memberList)
	http.HandleFunc("/depositList", depositList)
	http.HandleFunc("/depositEdit", depositEdit)
	http.HandleFunc("/depositUpdate", depositUpdate)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
