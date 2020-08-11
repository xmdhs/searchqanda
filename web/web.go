package web

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/yanyiwu/gojieba"
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
	if err != nil {
		e(w, err)
		return
	}
	if len(r) == 0 {
		x := gojieba.NewJieba(`dict/jieba.dict.utf8`, `dict/hmm_model.utf8`, `dict/user.dict.utf8`, `dict/idf.utf8`, `dict/stop_words.utf8`)
		defer x.Free()
		s := x.CutForSearch(query, true)
		t := strings.Join(s, " ")
		err := errors.New("未搜索到结果，建议使用\n\n" + t + "\n\n来尝试搜索")
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

func Index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(index))
}
