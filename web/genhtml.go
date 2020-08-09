package web

import (
	"html/template"
	"io"
)

func Genhtml() {

}

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
}

func pase(w io.Writer, list []resultslist, Name, page string) {
	T := true
	Link := ""
	if len(list) != 20 {
		T = false
	} else {
		Link = "/?q=" + Name + "&page=" + page
	}
	r := results{
		Name: Name,
		Link: Link,
		List: list,
		T:    T,
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
