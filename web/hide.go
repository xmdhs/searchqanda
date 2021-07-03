package web

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/xmdhs/searchqanda/get"
)

func Hidethead(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	var page string
	if len(q["page"]) == 0 {
		page = "0"
	} else {
		page = q["page"][0]

	}
	if len(q["q"]) != 0 {
		value := q["q"][0]
		i, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			e(w, err)
			return
		}
		offset := strconv.FormatInt(i*20, 10)
		i++
		page = strconv.FormatInt(i, 10)
		showhide(value, offset, page, w)
	} else {
		rows, err := get.Db.Query(`SELECT DISTINCT fid FROM hidethread`)
		if err != nil {
			e(w, err)
			return
		}
		defer rows.Close()
		var fid string
		list := make([]resultslist, 0)
		var rr resultslist
		rr.Title = "ALL"
		rr.Link = "./hide?q=all"
		list = append(list, rr)
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
		err = t.ExecuteTemplate(w, "html", r)
		if err != nil {
			panic(err)
		}
	}
}

func showhide(fid, offset, page string, w http.ResponseWriter) {
	var rows *sql.Rows
	var err error
	if fid != "all" {
		rows, err = get.Db.Query(`SELECT tid,dateline,authorid,author,subject FROM hidethread WHERE fid = ? ORDER BY tid DESC LIMIT 20 OFFSET ?`, fid, offset)
	} else {
		rows, err = get.Db.Query(`SELECT tid,dateline,authorid,author,subject FROM hidethread ORDER BY tid DESC LIMIT 20 OFFSET ?`, offset)
	}
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
	pase(w, list, fid, page, "./hide?q=")

}
