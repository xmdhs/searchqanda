package web

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
)

type results struct {
	Name string
	List []resultslist
	Link string
	T    bool
}

type resultslist struct {
	Title string
	Link  string
	Txt   string
	Txt1  string
	Key   string
}

func paseApi(w io.Writer, list []resultslist, Name, page, link string) {
	a := apiData{}
	a.Code = 0
	a.Data = list
	b, err := json.Marshal(a)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(b)
}

func pase(w io.Writer, list []resultslist, Name, page, link string) {
	T := true
	Link := ""
	if len(list) != 20 {
		T = false
	} else {
		Link = link + Name + "&page=" + page
	}
	r := results{
		Name: Name,
		Link: Link,
		List: list,
		T:    T,
	}
	err := t.ExecuteTemplate(w, "html", r)
	if err != nil {
		log.Println(err)
		return
	}
}

var t *template.Template

func init() {
	var err error
	t, err = template.ParseFS(htmlfs, "html/*")
	if err != nil {
		panic(err)
	}
}
