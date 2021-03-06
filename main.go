package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/mattn/go-sqlite3"

	bk "github.com/zzz404/MoneyGo/internal/bank"
	"github.com/zzz404/MoneyGo/internal/db"
	dp "github.com/zzz404/MoneyGo/internal/deposit"
	mb "github.com/zzz404/MoneyGo/internal/member"
	ut "github.com/zzz404/MoneyGo/internal/utils"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			err := db.DB.Close()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("db closed")
			}
			os.Exit(0)
		}
	}()

	startWeb()
}

func index(r *ut.HttpRequest, w *ut.HttpResponse) {
	tpl, err := ut.GetTemplate("/index.html")
	if w.ResponseForError(err) {
		return
	}
	all_map, all_total, time_map, time_total, err := dp.DepService.QueryTotalTWD()
	if w.ResponseForError(err) {
		return
	}
	var members []*mb.Member
	for _, member := range mb.Members {
		memberAllTotal, ok := all_map[member.Id]
		if ok {
			member.AllTotalTWD = memberAllTotal
		} else {
			member.AllTotalTWD = 0
		}
		memberTimeTotal, ok := time_map[member.Id]
		if ok {
			member.TimeTotalTWD = memberTimeTotal
		} else {
			member.TimeTotalTWD = 0
		}
		members = append(members, member)
	}

	err = tpl.Execute(w, map[string]interface{}{
		"members":            members,
		"allTotalTWDString":  fmt.Sprintf("%.2f", all_total),
		"timeTotalTWDString": fmt.Sprintf("%.2f", time_total),
	})

	w.ResponseForError(err)
}

func startWeb() {
	ut.HandleFunc("/", index)
	ut.HandleFunc("/accountList", bk.BankAccController.List)
	ut.HandleFunc("/accountEdit", bk.BankAccController.Edit)
	ut.HandleFunc("/accountAdd", bk.BankAccController.Add)
	ut.HandleFunc("/accountDelete", bk.BankAccController.Delete)
	ut.HandleFunc("/depositList", dp.DepController.List)
	ut.HandleFunc("/depositEdit", dp.DepController.Edit)
	ut.HandleFunc("/depositUpdate", dp.DepController.Update)
	ut.HandleFunc("/depositDelete", dp.DepController.Delete)
	ut.HandleFunc("/timeDepositList", dp.TimeDepController.List)

	fs := http.FileServer(http.Dir("Webapp/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("\nhttp://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error! ?????????????????????: %s\n", err.Error())
	}
}
