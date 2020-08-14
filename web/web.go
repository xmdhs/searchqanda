package web

import (
	"errors"
	"net/http"
	"strconv"
)

func WebRoot(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	var page, query string
	if len(q["q"]) == 0 {
		query = ""
	} else {
		query = q["q"][0]
	}
	if len(q["page"]) == 0 {
		page = "0"
	} else {
		page = q["page"][0]
	}
	if len(query) > 100 {
		err := errors.New("关键词过长")
		e(w, err)
		return
	}
	i, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		e(w, err)
		return
	}
	page = strconv.FormatInt(i*20, 10)
	r, err := search(query, page)
	if err != nil {
		e(w, err)
		return
	}
	if len(r) == 0 {
		http.NotFound(w, req)
		return
	}
	i++
	page = strconv.FormatInt(i, 10)
	pase(w, r, query, page)
}

func e(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}

func Style(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/css")
	w.Header().Set("Cache-Control", "max-age=315360000")
	w.Write([]byte(css))
}

func Index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(index))
}
