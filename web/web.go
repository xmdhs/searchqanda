package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	_ "embed"
)

func WebRoot(w http.ResponseWriter, req *http.Request) {
	b, err := htmlfs.ReadFile("html/web.html")
	if err != nil {
		panic(err)
	}
	w.Write([]byte(b))
}

func SerchApi(w http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")
	page := req.FormValue("page")
	if page == "" {
		page = "0"
	}
	if len(query) > 100 {
		apiErr("关键词过长", w)
		return
	}
	i, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		apiErr(err.Error(), w)
		return
	}
	page = strconv.FormatInt(i*20, 10)
	r, count, err := search(query, page)
	if err != nil {
		apiErr(err.Error(), w)
		return
	}
	if len(r) == 0 {
		apiErr("none", w)
		return
	}
	i++
	a := apiResults{
		List:  r,
		Count: count,
	}
	paseApi(w, a)
}

func apiErr(err string, w http.ResponseWriter) {
	e := apiData{
		Code: -1,
		Msg:  err,
	}
	w.WriteHeader(500)
	w.Write(e.Byte())
}

type apiData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (e *apiData) Byte() []byte {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return b
}

func e(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}

func Index(w http.ResponseWriter, req *http.Request) {
	b, err := htmlfs.ReadFile("html/index.html")
	if err != nil {
		panic(err)
	}
	w.Write([]byte(b))
}
