package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	s := req.Header.Get("sign")
	b, err := hex.DecodeString(s)
	if err != nil {
		e := apiData{
			Code: -1,
			Msg:  err.Error(),
		}
		w.WriteHeader(500)
		w.Write(e.Byte())
		return
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(req.URL.RawQuery))
	if !hmac.Equal(h.Sum(nil), b) {
		e := apiData{
			Code: -1,
			Msg:  "sign error",
		}
		w.WriteHeader(500)
		w.Write(e.Byte())
		return
	}
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

func Wasm(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/wasm")
	w.Header().Set("Cache-Control", "max-age=31536000")
	w.Write(wasm)
}

//go:embed s.wasm
var wasm []byte

var key = "12345678901234567890123456789012"
