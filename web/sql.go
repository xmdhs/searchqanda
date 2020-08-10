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
	sqline := strings.Builder{}
	sqline.WriteString(`SELECT tid,subject,txt FROM qa WHERE txt LIKE ?`)
	for i := 1; i < len(list); i++ {
		sqline.WriteString(` AND txt LIKE ?`)
	}
	sqline.WriteString(` ORDER BY tid DESC`)
	sqline.WriteString(` LIMIT 20 OFFSET ` + offset)
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(15*time.Second, func() {
		cancel()
	})
	l := make([]interface{}, 0, len(list))
	for _, v := range list {
		l = append(l, `%`+v+`%`)
	}
	rows, err := get.Db.QueryContext(ctx, sqline.String(), l...)
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
				if strings.Contains(strings.ToTitle(v.Message), strings.ToTitle(t)) {
					re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
					src := re.ReplaceAllString(v.Message, "")
					src = strings.ReplaceAll(src, "&nbsp;", "")
					a := strings.Index(strings.ToTitle(src), strings.ToTitle(t))
					aa := a - 150
					b := a + 150
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
					tt = strings.ToValidUTF8(tt, "")
					break
				}
			}
		}
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
