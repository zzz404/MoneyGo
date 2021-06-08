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

func getIntParameter(r *http.Request, name string, required bool) (int, bool, error) {
	strValue := r.URL.Query().Get(name)
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

func getFloatParameter(r *http.Request, name string, required bool) (float32, bool, error) {
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

func responseForError(err error, w http.ResponseWriter) bool {
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return true
	}
	return false
}
