package web

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/xmdhs/hidethread/get"
)

func search(txt, offset string) ([]resultslist, error) {
	if txt == "" {
		return []resultslist{}, errors.New(`txt == ""`)
	}
	list := strings.Split(txt, " ")
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(10*time.Second, func() {
		cancel()
	})
	txt = replace(txt)
	if txt == "" {
		return []resultslist{}, errors.New(`""`)
	}
	txt = "'" + txt + "'"
	rows, err := get.Db.QueryContext(ctx, `SELECT key,subject,source FROM qafts5 WHERE qafts5 MATCH `+txt+` ORDER BY rank DESC`)
	defer rows.Close()
	if err != nil {
		return []resultslist{}, err
	}
	var tid string
	var subject string
	var j string
	lists := make([]resultslist, 0, 20)
	for rows.Next() {
		err := rows.Scan(&tid, &subject, &j)
		if err != nil {
			return []resultslist{}, err
		}
		p := make([]post, 0)
		err = json.Unmarshal([]byte(j), &p)
		if err != nil {
			return []resultslist{}, err
		}
		var tt string
		for _, v := range p {
			for _, t := range list {
				t = strings.ReplaceAll(t, `"`, "")
				re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
				src := re.ReplaceAllString(v.Message, "")
				src = strings.ReplaceAll(src, "&nbsp;", "")
				src = strings.ReplaceAll(src, "/", "")
				if strings.Contains(strings.ToTitle(v.Message), strings.ToTitle(t)) {
					a := strings.Index(strings.ToTitle(src), strings.ToTitle(t))
					aa := a - 200
					b := a + 200
					if aa <= 0 {
						aa = 0
					}
					if b >= len(src) {
						b = len(src) - 1
					}
					if b == -1 {
						b = 0
					}
					if aa == -1 {
						aa = 0
					}
					tt = src[aa:b]
					break
				}
				if len(tt) == 0 {
					if len(src) <= 500 {
						tt = src
					} else {
						tt = src[0:500]
					}
				}
			}
		}
		tt = strings.ToValidUTF8(tt, "")
		l := resultslist{
			Title: subject,
			Link:  `https://www.mcbbs.net/thread-` + tid + `-1-1.html`,
			Txt:   tt,
		}
		lists = append(lists, l)
	}
	return lists, nil
}

type post struct {
	Message  string
	Authorid string
}

func replace(txt string) string {
	txt = strings.ReplaceAll(txt, ";", "")
	txt = strings.ReplaceAll(txt, "'", "")
	txt = strings.ReplaceAll(txt, ",", "")
	txt = strings.ReplaceAll(txt, "?", "")
	txt = strings.ReplaceAll(txt, "--", "")
	txt = strings.ReplaceAll(txt, "<", "")
	txt = strings.ReplaceAll(txt, ">", "")
	txt = strings.ReplaceAll(txt, "@", "")
	txt = strings.ReplaceAll(txt, "=", "")
	txt = strings.ReplaceAll(txt, "+", "")
	txt = strings.ReplaceAll(txt, "*", "")
	txt = strings.ReplaceAll(txt, "&", "")
	txt = strings.ReplaceAll(txt, "#", "")
	txt = strings.ReplaceAll(txt, "%", "")
	txt = strings.ReplaceAll(txt, "$", "")
	txt = strings.ReplaceAll(txt, `\`, "")
	txt = strings.ReplaceAll(txt, `(`, "")
	txt = strings.ReplaceAll(txt, `)`, "")
	txt = strings.ReplaceAll(txt, " -", "NOT ")
	txt = strings.ReplaceAll(txt, ".", "+")
	txt = strings.ReplaceAll(txt, "/", "+")

	return txt
}
