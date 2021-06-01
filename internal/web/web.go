package web

import (
	"html/template"
	"net/http"

	"github.com/zzz404/MoneyGo/internal/db"
)

func PersonsView(w http.ResponseWriter, r *http.Request) {
	members := [3]db.Member{{Id: 1, Name: "aaa"}, {Id: 2, Name: "bbb"}, {Id: 3, Name: "ccc"}}
	t, _ := template.ParseFiles("Webapp/PersonsView.html")
	t.Execute(w, map[string]interface{}{
		"persons": members,
	})
}
