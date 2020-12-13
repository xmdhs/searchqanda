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

const key = "c99be5fb4ebf56bf49ee9b7a7b5cfe4c6a6196a3ff4a8acdacb15dec27f37514"
