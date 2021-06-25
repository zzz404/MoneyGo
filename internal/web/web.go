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
	m, totalTWD, err := dp.QueryTotalTWD()
	if w.responseForError(err) {
		return
	}
	var members []*mb.Member
	for _, member1 := range mb.Members {
		member2 := member1
		memberTotal, ok := m[member2.Id]
		if ok {
			member2.TotalTWD = memberTotal
		} else {
			member2.TotalTWD = 0
		}
		members = append(members, member2)
	}

	err = tpl.Execute(w, map[string]interface{}{
		"members":        members,
		"totalTWDString": fmt.Sprintf("%.2f", totalTWD),
	})

	w.responseForError(err)
}

func readParameters(f *dp.QueryForm, r *HttpRequest) error {
	memberId, ok, err := r.getIntParameter("memberId", false)
	if err != nil {
		return err
	}
	if ok {
		f.MemberId = memberId
	}

	f.BankId, _, err = r.getIntParameter("bankId", false)
	if err != nil {
		return err
	}

	f.TypeCode, _, err = r.getIntParameter("typeCode", false)
	if err != nil {
		return err
	}

	f.CoinTypeCode, _, err = r.getParameter("coinTypeCode", false)
	if err != nil {
		return err
	}

	return nil
}

func depositList(r *HttpRequest, w *HttpResponse) {
	form := new(dp.QueryForm)
	err := readParameters(form, r)

	if w.responseForError(err) {
		return
	}
	tpl, err := getTemplate("/depositList.html")
	if w.responseForError(err) {
		return
	}
	deposits, err := dp.QueryDeposits(form)
	if w.responseForError(err) {
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
		"types":     dp.DepositTypes,
		"deposits":  deposits,
		"count":     len(deposits),
		"totalTWD":  fmt.Sprintf("%.2f", totalTWD),
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
		memberId, _, err := r.getIntParameter("memberId", false)
		if w.responseForError(err) {
			return
		}
		bankId, _, err := r.getIntParameter("bankId", false)
		if w.responseForError(err) {
			return
		}
		typeCode, _, err := r.getIntParameter("typeCode", false)
		if w.responseForError(err) {
			return
		}
		coinTypeCode, _, err := r.getParameter("coinTypeCode", false)
		if w.responseForError(err) {
			return
		}
		deposit = &dp.Deposit{MemberId: memberId, BankId: bankId, TypeCode: typeCode, CoinTypeCode: coinTypeCode}
	}

	data := map[string]interface{}{
		"deposit":      deposit,
		"members":      mb.Members,
		"banks":        bk.Banks,
		"bankAccounts": bk.BankIdAccountsMapJson,
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

	deposit.MemberId, _, err = r.getIntParameter("memberId", true)
	if w.responseForError(err) {
		return
	}

	deposit.BankId, _, err = r.getIntParameter("bankId", true)
	if w.responseForError(err) {
		return
	}

	deposit.TypeCode, _, err = r.getIntParameter("typeCode", true)
	if w.responseForError(err) {
		return
	}

	deposit.Amount, _, err = r.getFloatParameter("amount", true)
	if w.responseForError(err) {
		return
	}

	deposit.CoinTypeCode, _, err = r.getParameter("coinTypeCode", true)
	if w.responseForError(err) {
		return
	}

	if hasId {
		err = dp.UpdateDeposit(deposit)
		if w.responseForError(err) {
			return
		}
	} else {
		_, err = dp.AddDeposit(deposit)
		if w.responseForError(err) {
			return
		}
	}
	w.Redirect("/static/ReloadOpenerThenClose.html", r)
}

func depositDelete(r *HttpRequest, w *HttpResponse) {
	id, _, err := r.getIntParameter("id", true)
	if w.responseJsonError(err) {
		return
	}

	err = dp.DeleteDeposit(id)
	if w.responseJsonError(err) {
		return
	}
	w.writeJson(true, "", nil)
}

func Start() {
	handleFunc("/", memberList)
	handleFunc("/depositList", depositList)
	handleFunc("/depositEdit", depositEdit)
	handleFunc("/depositUpdate", depositUpdate)
	handleFunc("/depositDelete", depositDelete)
	fs := http.FileServer(http.Dir("Webapp/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("\nhttp://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error! 結束時發生錯誤: %s\n", err.Error())
	}
}
