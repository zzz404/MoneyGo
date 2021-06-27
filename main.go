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
	m, totalTWD, err := dp.DepService.QueryTotalTWD()
	if w.ResponseForError(err) {
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
	fs := http.FileServer(http.Dir("Webapp/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("\nhttp://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error! 結束時發生錯誤: %s\n", err.Error())
	}
}
