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
		r.HandleFunc("/search/hide", web.Auth(web.Hidethead, key))
		r.HandleFunc("/search/snapshot", web.Auth(web.Snapshot, key))
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

const key = "53e6b64604b8ed484bb6c67b93c0987bec828db8d4a725d080dd7092b9fb15b2"
