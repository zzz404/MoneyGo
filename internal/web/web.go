package web

import (
	"fmt"
	"net/http"

	bk "github.com/zzz404/MoneyGo/internal/bank"
	"github.com/zzz404/MoneyGo/internal/coin"
	dp "github.com/zzz404/MoneyGo/internal/deposit"
	mb "github.com/zzz404/MoneyGo/internal/member"
)

func memberList(r *HttpRequest, w *HttpResponse) {
	tpl, err := getTemplate("/memberList.html")
	if w.responseForError(err) {
		return
	}
	members := mb.Members
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
	id, isEdit, err := r.getIntParameter("id", false)
	if w.responseForError(err) {
		return
	}
	var deposit *dp.Deposit
	if isEdit {
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
		member, err := mb.GetMember(memberId)
		if w.responseForError(err) {
			return
		}
		deposit = &dp.Deposit{Member: member}
	}

	data := map[string]interface{}{
		"deposit":      deposit,
		"banks":        bk.Banks,
		"depositTypes": dp.DepositTypes,
		"coinTypes":    coin.CoinTypes,
	}
	if isEdit {
		data["id"] = deposit.Id
	}
	err = tpl.Execute(w, data)

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

	memberId, _, err := r.getIntParameter("memberId", true)
	if w.responseForError(err) {
		return
	}
	member, err := mb.GetMember(memberId)
	if w.responseForError(err) {
		return
	}
	deposit.Member = member

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

	coinTypeCode, _, err := r.getParameter("coinTypeCode", true)
	if w.responseForError(err) {
		return
	}

	deposit.CoinType, err = coin.GetCoinTypeByCode(coinTypeCode)
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
	w.Redirect(url, r)
}

func depositDelete(r *HttpRequest, w *HttpResponse) {
	id, _, err := r.getIntParameter("id", true)
	if w.responseForError(err) {
		return
	}
	err = dp.DeleteDeposit(id)
	if w.responseForError(err) {
		return
	}
	w.Redirect("/depositList", r)
}

func Start() {
	handleFunc("/", memberList)
	handleFunc("/depositList", depositList)
	handleFunc("/depositEdit", depositEdit)
	handleFunc("/depositUpdate", depositUpdate)
	handleFunc("/depositDelete", depositDelete)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error! 結束時發生錯誤: %s\n", err.Error())
	}
}
