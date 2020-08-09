package main

import (
	"log"
	"net/http"
	"time"

	"github.com/xmdhs/hidethread/get"
	"github.com/xmdhs/hidethread/web"
)

func main() {
	go upsql()
	r := http.NewServeMux()
	r.HandleFunc("/s", web.WebRoot)
	r.HandleFunc("/style.css", web.Style)
	r.HandleFunc("/hide", web.Hidethead)
	s := http.Server{
		Addr:         ":8081",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 12 * time.Second,
		Handler:      r,
	}
	log.Println(s.ListenAndServe())
}

func upsql() {
	for {
		get.Startrange()
		time.Sleep(24 * time.Hour)
	}
}
