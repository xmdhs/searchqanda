package web

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/xmdhs/searchqanda/get"
)

func search(txt, offset string) ([]resultslist, error) {
	if txt == "" {
		return []resultslist{}, errors.New(`""`)
	}

	if txt == "" {
		return []resultslist{}, errors.New(`txt == ""`)
	}
	list := cut(txt)
	l := strings.Split(txt, " ")
	for i, v := range list {
		v = replace(v)
		list[i] = cutsearch(v)
	}
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(10*time.Second, func() {
		cancel()
	})
	txt = strings.Join(list, " ")
	txt = "'" + txt + "'"
	rows, err := get.Db.QueryContext(ctx, `SELECT key,subject,source FROM qafts5 WHERE qafts5 MATCH `+txt+` ORDER BY rank DESC LIMIT 20 OFFSET ?`, offset)
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
			src := strings.ReplaceAll(v.Message, "/", "")
			for _, t := range l {
				t = strings.ReplaceAll(t, `"`, "")
				if strings.Contains(strings.ToTitle(src), strings.ToTitle(t)) {
					re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
					src = re.ReplaceAllString(src, "")
					src = strings.ReplaceAll(src, "&nbsp;", "")
					src = strings.ReplaceAll(src, "\n;", "")
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
					tt = strings.ReplaceAll(tt, t, `<b>`+t+`</b>`)
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
			Txt:   template.HTML(tt),
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
	txt = strings.ReplaceAll(txt, "+", " ")
	txt = strings.ReplaceAll(txt, ";", " ")
	txt = strings.ReplaceAll(txt, "'", " ")
	txt = strings.ReplaceAll(txt, ",", " ")
	txt = strings.ReplaceAll(txt, "?", " ")
	txt = strings.ReplaceAll(txt, "--", " ")
	txt = strings.ReplaceAll(txt, "<", " ")
	txt = strings.ReplaceAll(txt, ">", " ")
	txt = strings.ReplaceAll(txt, "@", " ")
	txt = strings.ReplaceAll(txt, "=", " ")
	txt = strings.ReplaceAll(txt, "*", " ")
	txt = strings.ReplaceAll(txt, ":", " ")
	txt = strings.ReplaceAll(txt, "&", " ")
	txt = strings.ReplaceAll(txt, "#", " ")
	txt = strings.ReplaceAll(txt, "%", " ")
	txt = strings.ReplaceAll(txt, "$", " ")
	txt = strings.ReplaceAll(txt, `\`, " ")
	txt = strings.ReplaceAll(txt, `(`, " ")
	txt = strings.ReplaceAll(txt, `)`, " ")
	txt = strings.ReplaceAll(txt, ".", " ")
	txt = strings.ReplaceAll(txt, "/", " ")

	return txt
}

func cut(txt string) []string {
	ss := make([]string, 0)
	txt = txt + " "
	s := strings.Builder{}
	t := true
	for _, v := range txt {
		if v == 45 {
			s.WriteRune(v)
			continue
		}
		if v == 34 {
			if t {
				t = false
			} else {
				t = true
			}
		}
		s.WriteRune(v)
		if v == 32 && t {
			t := strings.Trim(s.String(), " ")
			ss = append(ss, t)
			s.Reset()
		}
	}
	if len(ss) == 0 {
		ss = append(ss, txt)
	}
	return ss
}

func cutsearch(src string) string {
	remove := false
	if strings.HasPrefix(src, "-") {
		src = strings.TrimPrefix(src, `-`)
		remove = true
	}
	src = strings.ReplaceAll(src, `"`, "")
	src = strings.ReplaceAll(src, "-", " ")
	l := get.Seg.CutSearch(src, true)
	ll := make([]string, 0)

	for _, v := range l {
		v = strings.Trim(v, " ")
		if v != "" {
			ll = append(ll, v)
		}
	}
	src = strings.Join(ll, " ")
	if remove {
		return `NOT "` + src + `"`
	}
	return `"` + src + `"`
}
