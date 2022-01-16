package utils

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

func GetTemplate(path string) (*template.Template, error) {
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

func (w *HttpResponse) ResponseForError(err error) bool {
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return true
	}
	return false
}

func (res *HttpResponse) Redirect(url string, req *HttpRequest) {
	http.Redirect(res, req.Request, url, http.StatusSeeOther)
}

func (r *HttpRequest) GetParameter(name string, required bool) (string, bool, error) {
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

func (r *HttpRequest) GetIntParameter(name string, required bool) (int, bool, error) {
	strValue, found, err := r.GetParameter(name, required)
	if err != nil || !found {
		return 0, found, err
	} else {
		value, err := strconv.Atoi(strValue)
		return value, found, err
	}
}

func (r *HttpRequest) GetIntPointerParameter(name string, required bool) (*int, error) {
	strValue, found, err := r.GetParameter(name, required)
	if err != nil || !found {
		return nil, err
	} else {
		value, err := strconv.Atoi(strValue)
		return &value, err
	}
}

func (r *HttpRequest) GetFloatParameter(name string, required bool) (float64, bool, error) {
	strValue, found, err := r.GetParameter(name, required)
	if err != nil || !found {
		return 0, found, err
	} else {
		value, err := strconv.ParseFloat(strValue, 32)
		return value, found, err
	}
}

func (r *HttpRequest) GetFloatPointerParameter(name string, required bool) (*float64, error) {
	strValue, found, err := r.GetParameter(name, required)
	if err != nil || !found {
		return nil, err
	} else {
		value, err := strconv.ParseFloat(strValue, 32)
		return &value, err
	}
}

func (r *HttpRequest) GetBoolParameter(name string, required bool) (value bool, found bool, err error) {
	strValue, found, err := r.GetParameter(name, required)
	if err != nil || !found {
		return
	} else {
		if strValue == "true" {
			value = true
		} else if strValue == "false" {
			value = false
		} else {
			err = fmt.Errorf("字串 %s 轉 boolean 失敗", strValue)
		}
		return
	}
}

func (r *HttpRequest) GetBoolPointerParameter(name string, required bool) (*bool, error) {
	strValue, found, err := r.GetParameter(name, required)
	if err != nil || !found {
		return nil, err
	} else {
		if strValue == "true" {
			result := true
			return &result, nil
		} else if strValue == "false" {
			result := false
			return &result, nil
		} else {
			err = fmt.Errorf("字串 %s 轉 boolean 失敗", strValue)
			return nil, err
		}
	}
}

func (r *HttpRequest) GetDatePointerParameter(name string, required bool) (*time.Time, error) {
	strValue, found, err := r.GetParameter(name, required)
	if err != nil || !found {
		return nil, err
	} else {
		value, err := ParseDate(strValue)
		return value, err
	}
}

func HandleFunc(path string, fn func(r *HttpRequest, w *HttpResponse)) {
	f := func(w http.ResponseWriter, r *http.Request) {
		ww := &HttpResponse{w}
		rr := &HttpRequest{r, r.Method == "POST"}
		if rr.isPost {
			err := r.ParseForm()
			if err != nil {
				ww.ResponseForError(err)
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

func (w *HttpResponse) WriteJson(success bool, message string, data interface{}) {
	jsonResult := new(JsonResult)
	jsonResult.Success = success
	jsonResult.Data = data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResult)
}

func (w *HttpResponse) ResponseJsonError(err error) bool {
	if err != nil {
		w.WriteJson(false, err.Error(), nil)
		return true
	}
	return false
}

type Radio struct {
	Value   string
	Text    string
	Checked bool
}
