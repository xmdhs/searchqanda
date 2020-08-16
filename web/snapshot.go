package web

import (
	"encoding/json"
	"net/http"

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
		w.Write([]byte(snapshot))
	} else {
		rows, err := get.Db.Query(`SELECT source FROM qafts5 WHERE key = ?`, query)
		if err != nil {
			e(w, err)
			return
		}
		source := ""
		rows.Next()
		rows.Scan(&source)
		rows.Close()
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
			}
			rlist = append(rlist, r)
		}
		pase(w, rlist, query, "", "")
	}
}

type posts struct {
	Message  string
	Authorid string
}
