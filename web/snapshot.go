package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/xmdhs/searchqanda/get"
)

func Snapshot(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var query string
	if len(q["q"]) == 0 {
		query = ""
	} else {
		query = q["q"][0]
	}
	tojson := false
	if len(q["json"]) != 0 {
		tojson = true
	}
	if query == "" {
		b, err := htmlfs.ReadFile("html/snapshot.html")
		if err != nil {
			panic(err)
		}
		w.Write([]byte(b))
	} else {
		i, err := strconv.ParseInt(query, 10, 64)
		if err != nil {
			e(w, err)
			return
		}
		s := strconv.FormatInt(i, 10)
		rows := get.Db.QueryRow(`SELECT source FROM qafts5 WHERE key MATCH ?`, s)
		source := ""
		err = rows.Scan(&source)
		if err != nil {
			e(w, err)
			return
		}
		if source == "" {
			http.NotFound(w, r)
			return
		}
		if tojson {
			w.Write([]byte(source))
			return
		}
		j := make([]posts, 0)
		err = json.Unmarshal([]byte(source), &j)
		if err != nil {
			e(w, err)
			return
		}
		rlist := make([]resultslist, 0)
		for _, v := range j {
			r := resultslist{
				Title: v.Authorid,
				Txt:   v.Message,
				Link:  "https://www.mcbbs.net/?" + v.Authorid,
			}
			rlist = append(rlist, r)
		}
		pase(w, rlist, "快照 - "+query, "", "")
	}
}

type posts struct {
	Message  string
	Authorid string
}
