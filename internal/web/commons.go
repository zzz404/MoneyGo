package web

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"
)

type myTpl struct {
	tpl          *template.Template
	err          error
	modifiedtime time.Time
}

var templateMap map[string]*myTpl = make(map[string]*myTpl)

func getTemplate(path string) (*template.Template, error) {
	realPath := "Webapp" + path
	file, err := os.Stat(realPath)
	if err != nil {
		return nil, err
	}
	modifiedtime := file.ModTime()

	var reload bool
	tpl, found := templateMap[path]
	if found {
		reload = modifiedtime.After(tpl.modifiedtime)
	} else {
		reload = true
	}
	if reload {
		t, err := template.ParseFiles(realPath)
		tpl = &myTpl{tpl: t, err: err, modifiedtime: modifiedtime}
		templateMap[path] = tpl
	}
	return tpl.tpl, tpl.err
}

func lackParamError(name string) error {
	return fmt.Errorf("缺少必要參數 : %s", name)
}

type HttpRequest struct {
	*http.Request
	isPost bool
}

type HttpResponse struct {
	http.ResponseWriter
}

func (w HttpResponse) responseForError(err error) bool {
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return true
	}
	w.Header()
	return false
}

func (r *HttpRequest) getIntParameter(name string, required bool) (int, bool, error) {
	var strValue string
	if r.isPost {
		strValue = r.FormValue(name)
	} else {
		strValue = r.URL.Query().Get(name)
	}
	if strValue == "" {
		if required {
			return 0, false, lackParamError(name)
		} else {
			return 0, false, nil
		}
	} else {
		i, err := strconv.Atoi(strValue)
		return i, true, err
	}
}

func (r *HttpRequest) getFloatParameter(name string, required bool) (float32, bool, error) {
	strValue := r.URL.Query().Get(name)
	if strValue == "" {
		if required {
			return 0, false, lackParamError(name)
		} else {
			return 0, false, nil
		}
	}
	value, err := strconv.ParseFloat(strValue, 32)
	if err != nil {
		return 0, true, err
	} else {
		return float32(value), true, nil
	}
}

func handleFunc(path string, fn func(r *HttpRequest, w *HttpResponse)) {
	f := func(w http.ResponseWriter, r *http.Request) {
		ww := &HttpResponse{w}
		rr := &HttpRequest{r, r.Method == "POST"}
		if rr.isPost {
			err := r.ParseForm()
			if err != nil {
				ww.responseForError(err)
				return
			}
		}
		fn(rr, ww)
	}
	http.HandleFunc(path, f)
}
