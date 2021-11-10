package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/xmdhs/searchqanda/get"
)

var htmlreg = regexp.MustCompile(`\<[\S\s]+?\>`)

var ErrEmpty = errors.New("empty")

func search(txt, offset string) ([]resultslist, int, error) {
	if txt == "" {
		return []resultslist{}, 0, fmt.Errorf("search: %w", ErrEmpty)
	}
	list := cut(txt)
	l := strings.Split(txt, " ")
	for i, v := range list {
		list[i] = cutsearch(v)
	}
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(10*time.Second, func() {
		cancel()
	})
	txt = strings.Join(list, " ")
	rows, err := get.Db.QueryContext(ctx, `SELECT key,subject,source FROM qafts5 WHERE qafts5 MATCH ? ORDER BY rank LIMIT 20 OFFSET ?`, txt, offset)
	if err != nil {
		return []resultslist{}, 0, fmt.Errorf("search: %w", err)
	}
	defer rows.Close()
	var tid string
	var subject string
	var j string
	lists := make([]resultslist, 0, 20)
	for rows.Next() {
		err := rows.Scan(&tid, &subject, &j)
		if err != nil {
			return []resultslist{}, 0, fmt.Errorf("search: %w", err)
		}
		p := make([]post, 0)
		err = json.Unmarshal([]byte(j), &p)
		if err != nil {
			return []resultslist{}, 0, fmt.Errorf("search: %w", err)
		}
		var b1, b2, key string
		for _, v := range p {
			src := strings.ReplaceAll(v.Message, "/", "")
			for _, t := range l {
				t = strings.ReplaceAll(t, `"`, "")
				t := strings.ReplaceAll(t, "/", "")
				if strings.Contains(strings.ToTitle(src), strings.ToTitle(t)) {
					src = htmlreg.ReplaceAllString(src, "")
					src = html.UnescapeString(html.UnescapeString(src))
					src = strings.ReplaceAll(src, "\n;", "")
					a := strings.Index(strings.ToTitle(src), strings.ToTitle(t))
					if a == -1 {
						continue
					}
					a1 := a + len(t)
					key = src[a:a1]
					aa := a - 200
					b := a1 + 200
					if aa <= 0 {
						aa = 0
					}
					if b > len(src) {
						b = len(src)
					}
					if b == -1 {
						b = 0
					}
					if aa == -1 {
						aa = 0
					}
					b1 = src[aa:a]
					b2 = src[a1:b]
					break
				}
				if len(b1) == 0 {
					if len(src) <= 500 {
						b1 = src
					} else {
						b1 = src[0:500]
					}
				}
			}
		}
		b1 = strings.ToValidUTF8(b1, "")
		b2 = strings.ToValidUTF8(b2, "")
		l := resultslist{
			Title: subject,
			Link:  `https://www.mcbbs.net/thread-` + tid + `-1-1.html`,
			Txt:   b1,
			Txt1:  b2,
			Key:   key,
		}
		lists = append(lists, l)
	}
	count := 0
	row := get.Db.QueryRowContext(ctx, `SELECT COUNT(*) FROM qafts5 WHERE qafts5 MATCH ?`, txt)
	err = row.Scan(&count)
	if err != nil {
		return []resultslist{}, 0, fmt.Errorf("search: %w", err)
	}
	return lists, count, nil
}

type post struct {
	Message  string
	Authorid string
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
	l := get.X.CutForSearch(src, true)
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
