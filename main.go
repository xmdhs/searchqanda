package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/xmdhs/searchqanda/get"
	"github.com/xmdhs/searchqanda/web"
)

func main() {
	if len(os.Args) != 1 {
		upsql()
	} else {
		r := http.NewServeMux()
		r.HandleFunc("/search", web.Index)
		r.HandleFunc("/search/s", web.WebRoot)
		r.HandleFunc("/search/style.css", web.Style)
		r.HandleFunc("/search/hide", web.Hidethead)
		s := http.Server{
			Addr:         ":8081",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 20 * time.Second,
			Handler:      r,
		}
		log.Println(s.ListenAndServe())
	}
}

func upsql() {
	for i := 0; i < 2; i++ {
		get.Startrange()
	}
}
