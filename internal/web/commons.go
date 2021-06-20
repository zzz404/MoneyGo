package web

import (
	"encoding/json"
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

func (w *HttpResponse) responseForError(err error) bool {
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return true
	}
	return false
}

func (res *HttpResponse) Redirect(url string, req *HttpRequest) {
	http.Redirect(res, req.Request, url, http.StatusSeeOther)
}

func (r *HttpRequest) getParameter(name string, required bool) (string, bool, error) {
	var strValue string
	if r.isPost {
		strValue = r.FormValue(name)
	} else {
		strValue = r.URL.Query().Get(name)
	}
	if strValue == "" {
		var err error
		if required {
			err = lackParamError(name)
		}
		return "", false, err
	}
	return strValue, true, nil
}

func (r *HttpRequest) getIntParameter(name string, required bool) (int, bool, error) {
	strValue, found, err := r.getParameter(name, required)
	if err != nil || !found {
		return 0, found, err
	} else {
		value, err := strconv.Atoi(strValue)
		return value, found, err
	}
}

func (r *HttpRequest) getFloatParameter(name string, required bool) (float64, bool, error) {
	strValue, found, err := r.getParameter(name, required)
	if err != nil || !found {
		return 0, found, err
	} else {
		value, err := strconv.ParseFloat(strValue, 32)
		return value, found, err
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

type JsonResult struct {
	Success bool
	Message string
	Data    interface{}
}

func (w *HttpResponse) writeJson(success bool, message string, data interface{}) {
	jsonResult := new(JsonResult)
	jsonResult.Success = success
	jsonResult.Data = data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResult)
}

func (w *HttpResponse) responseJsonError(err error) bool {
	if err != nil {
		w.writeJson(false, err.Error(), nil)
		return true
	}
	return false
}
