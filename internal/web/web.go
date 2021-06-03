package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/zzz404/MoneyGo/internal/db"
)

func PersonsView(w http.ResponseWriter, r *http.Request) {
	members, err := db.QueryMembers()
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
	} else {
		t, _ := template.ParseFiles("Webapp/PersonsView.html")
		t.Execute(w, map[string]interface{}{
			"persons": members,
		})
	}
}

func Start() {
	http.HandleFunc("/", PersonsView)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
