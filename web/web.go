package web

import (
	"encoding/json"
	"errors"
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
		err := errors.New("关键词过长")
		e := apiData{
			Code: -1,
			Msg:  err.Error(),
		}
		w.WriteHeader(500)
		w.Write(e.Byte())
		return
	}
	i, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		e := apiData{
			Code: -1,
			Msg:  err.Error(),
		}
		w.WriteHeader(500)
		w.Write(e.Byte())
		return
	}
	page = strconv.FormatInt(i*20, 10)
	r, err := search(query, page)
	if err != nil {
		e := apiData{
			Code: -1,
			Msg:  err.Error(),
		}
		w.WriteHeader(500)
		w.Write(e.Byte())
		return
	}
	if len(r) == 0 {
		e := apiData{
			Code: -1,
			Msg:  "none",
		}
		w.WriteHeader(404)
		w.Write(e.Byte())
		return
	}
	i++
	page = strconv.FormatInt(i, 10)
	paseApi(w, r, query, page, "./s?q=")
}

type apiData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data interface{}
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
