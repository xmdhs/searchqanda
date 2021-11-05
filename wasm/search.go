package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"syscall/js"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/prop"
	router "marwan.io/vecty-router"
)

type search struct {
	vecty.Core
}

func (s *search) Render() vecty.ComponentOrHTML {
	var c vecty.ComponentOrHTML
	href := js.Global().Get("location").Get("href").String()
	u, err := url.Parse(href)
	var query string
	var nextLink vecty.ComponentOrHTML
	if err != nil {
		c = elem.Paragraph(vecty.Text(err.Error()))
	} else {
		query := u.Query().Get("q")
		page := u.Query().Get("page")
		b := false
		c, b, err = body(query, page)
		if err != nil {
			c = elem.Paragraph(vecty.Text(err.Error()))
		}
		vecty.SetTitle(query + " - 问答版搜索")
		if b {
			q := url.Values{}
			q.Set("q", query)
			q.Set("page", page)
			nextLink = elem.Paragraph(router.Link("/search/s?"+q.Encode(), "next", router.LinkOptions{}))
		}
	}
	return elem.Body(
		elem.Div(
			vecty.Markup(
				vecty.Class("container-lg", "px-3", "my-5", "markdown-body"),
			),
			elem.Heading1(
				vecty.Text(query),
			),
			c,
			elem.HorizontalRule(),
			vecty.If(nextLink != nil, nextLink),
		),
	)
}

var (
	ErrTooLong  = errors.New("query too long")
	ErrNotFound = errors.New("not found")
)

func body(query, page string) (vecty.ComponentOrHTML, bool, error) {
	if page == "" {
		page = "0"
	}
	if len(query) > 100 {
		return nil, false, fmt.Errorf("body: %w", ErrTooLong)
	}
	i, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return nil, false, fmt.Errorf("body: %w", err)
	}

	page = strconv.FormatInt(i*20, 10)
	r, err := serachApi(query, page)
	if err != nil {
		return nil, false, fmt.Errorf("body: %w", err)
	}
	if len(r) == 0 {
		return nil, false, fmt.Errorf("body: %w", ErrNotFound)
	}
	b := false
	if len(r) == 20 {
		b = true
	}
	return renderList(r), b, nil
}

func renderList(list []resultslist) vecty.ComponentOrHTML {
	l := make(vecty.List, 0, len(list))
	for _, v := range list {
		l = append(l,
			elem.Heading3(
				elem.Anchor(
					vecty.Markup(
						prop.Href(v.Link),
						vecty.Attribute("target", "_blank"),
					),
					vecty.Text(v.Title),
				),
			),
			elem.BlockQuote(
				elem.Paragraph(
					vecty.Text(v.Txt),
					vecty.If(v.Key != "",
						vecty.Tag(
							"font",
							vecty.Markup(
								vecty.Attribute("color", "red"),
							),
							vecty.Text(v.Key),
						),
					),
					vecty.Text(v.Txt1),
				),
			),
		)
	}
	return l
}

var c = http.Client{
	Timeout: time.Second * 10,
}

var host string

func init() {
	host = js.Global().Get("location").Get("origin").String()
}

var key = "12345678901234567890123456789012"

func serachApi(txt string, offset string) ([]resultslist, error) {
	q := url.Values{}
	q.Set("q", txt)
	q.Set("page", offset)
	qu := q.Encode()
	d := dohmac([]byte(qu), []byte(key))
	req, err := http.NewRequest("GET", host+"/search/api/s?"+qu, nil)
	if err != nil {
		return nil, fmt.Errorf("serachApi: %w", err)
	}
	req.Header.Set("sign", d)
	reps, err := c.Do(req)
	if reps != nil {
		defer reps.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("serachApi: %w", err)
	}
	var r apiData
	err = json.NewDecoder(reps.Body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("serachApi: %w", err)
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("serachApi: %w: %v", ErrNotFound, r.Msg)
	}
	rl := []resultslist{}
	err = json.Unmarshal([]byte(r.Data), &rl)
	if err != nil {
		return nil, fmt.Errorf("serachApi: %w", err)
	}
	return rl, nil
}

type resultslist struct {
	Title string
	Link  string
	Txt   string
	Txt1  string
	Key   string
}

type apiData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data json.RawMessage
}

func dohmac(msg, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(msg)
	return hex.EncodeToString(h.Sum(nil))
}
