package web

import (
	"html/template"
	"net/http"
)

type Person struct {
	Id   int8
	Name string
}

func PersonsView(w http.ResponseWriter, r *http.Request) {
	persons := [3]Person{{Id: 1, Name: "aaa"}, {Id: 2, Name: "bbb"}, {Id: 3, Name: "ccc"}}
	t, _ := template.ParseFiles("Webapp/PersonsView.html")
	t.Execute(w, map[string]interface{}{
		"persons": persons,
	})
}
