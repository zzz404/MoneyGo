package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/zzz404/MoneyGo/internal/db"
)

func memberList(w http.ResponseWriter, r *http.Request) {
	members, err := db.QueryMembers()
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}
	t, err := template.ParseFiles("Webapp/memberList.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}
	t.Execute(w, map[string]interface{}{
		"members": members,
	})

}

func depositList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "zzz")
}

func Start() {
	http.HandleFunc("/", memberList)
	http.HandleFunc("/member/", depositList)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
