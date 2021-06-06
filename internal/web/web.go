package web

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/zzz404/MoneyGo/internal/db"
	dp "github.com/zzz404/MoneyGo/internal/deposit"
)

func memberList(w http.ResponseWriter, r *http.Request) {
	tpl, err := getTemplate("/memberList.html")
	if responseError(err, w) {
		return
	}
	members, err := db.QueryMembers()
	if responseError(err, w) {
		return
	}
	err = tpl.Execute(w, map[string]interface{}{
		"members": members,
	})
	responseError(err, w)
}

func depositList(w http.ResponseWriter, r *http.Request) {
	memberId, err := getIntParameter(r, "memberId")
	if responseError(err, w) {
		return
	}
	tpl, err := getTemplate("/depositList.html")
	if responseError(err, w) {
		return
	}
	deposits, err := dp.QueryDeposits(memberId)
	if responseError(err, w) {
		return
	}
	err = tpl.Execute(w, map[string]interface{}{
		"memberId": memberId,
		"deposits": deposits,
	})
	responseError(err, w)
}

func depositEdit(w http.ResponseWriter, r *http.Request) {
	tpl, err := getTemplate("/depositEdit.html")
	if responseError(err, w) {
		return
	}
	id, err := getIntParameter(r, "id")
	if responseError(err, w) {
		return
	}
	var deposit *dp.Deposit
	if id != 0 {
		deposit, err = dp.GetDeposit(id)
		if err == nil && deposit == nil {
			err = errors.New(fmt.Sprintf("Deposit %d 不存在!", id))
		}
		if responseError(err, w) {
			return
		}
	} else {
		deposit = &dp.Deposit{}
	}
	err = tpl.Execute(w, deposit)
	responseError(err, w)
}

func Start() {
	http.HandleFunc("/", memberList)
	http.HandleFunc("/depositList", depositList)
	http.HandleFunc("/depositEdit", depositEdit)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
