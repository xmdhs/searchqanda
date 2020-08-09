package web

import (
	"io"
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
	i, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		e(w, err)
		return
	}
	page = strconv.FormatInt(i*20, 10)
	r, err := search(query, page)
	if len(r) == 0 {
		http.NotFound(w, req)
		return
	}
	if err != nil {
		e(w, err)
		return
	}
	i++
	page = strconv.FormatInt(i, 10)
	pase(w, r, query, page)
}

func e(w io.Writer, err error) {
	w.Write([]byte(err.Error()))
}

func Style(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/css")
	w.Write([]byte(css))
}
