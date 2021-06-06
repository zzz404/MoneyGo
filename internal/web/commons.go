package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type myTpl struct {
	tpl *template.Template
	err error
}

var templateMap map[string]*myTpl = make(map[string]*myTpl)

func getTemplate(path string) (*template.Template, error) {
	tpl, ok := templateMap[path]
	if !ok || 1+1 == 2 {
		t, err := template.ParseFiles("Webapp" + path)
		tpl = &myTpl{tpl: t, err: err}
		templateMap[path] = tpl
	}
	return tpl.tpl, tpl.err
}

func getIntParameter(r *http.Request, name string) (int, error) {
	strValue := r.URL.Query().Get(name)
	return strconv.Atoi(strValue)
}

func responseError(err error, w http.ResponseWriter) bool {
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return true
	}
	return false
}
