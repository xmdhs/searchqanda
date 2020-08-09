package web

import (
	"html/template"
	"io"
	"net/http"

	"github.com/xmdhs/hidethread/get"
)

func Hidethead(w http.ResponseWriter, req *http.Request) {
	get.M.RLock()
	defer get.M.RUnlock()
	q := req.URL.Query()
	if len(q["q"]) != 0 {
		value := q["q"][0]
		showhide(value, w)
	} else {
		rows, err := get.Db.Query(`SELECT DISTINCT fid FROM hidethread`)
		defer rows.Close()
		if err != nil {
			e(w, err)
			return
		}
		var fid string
		list := make([]resultslist, 0)
		for rows.Next() {
			rows.Scan(&fid)
			var r resultslist
			r.Title = fid
			r.Link = "./hide?q=" + fid
			list = append(list, r)
		}
		r := results{
			Name: "无权查看的帖子",
			List: list,
		}
		t, err := template.New("page").Parse(html)
		if err != nil {
			panic(err)
		}
		err = t.Execute(w, r)
		if err != nil {
			panic(err)
		}
	}
}

func showhide(fid string, w io.Writer) {
	rows, err := get.Db.Query(`SELECT tid,dateline,authorid,author,subject FROM hidethread WHERE fid = ? ORDER BY tid DESC`, fid)
	defer rows.Close()
	if err != nil {
		e(w, err)
		return
	}
	list := make([]resultslist, 0)
	var tid, dateline, authorid, author, subject string
	for rows.Next() {
		rows.Scan(&tid, &dateline, &authorid, &author, &subject)
		var r resultslist
		r.Title = subject
		r.Link = "https://www.mcbbs.net/thread-" + tid + "-1-1.html"
		r.Txt = author + "(" + authorid + ")" + "  ---" + dateline
		list = append(list, r)
	}
	r := results{
		Name: fid,
		List: list,
	}
	t, err := template.New("page").Parse(html)
	if err != nil {
		panic(err)
	}
	err = t.Execute(w, r)
	if err != nil {
		panic(err)
	}

}
