package main

import (
	"log"
	"net/http"
	"time"

	"github.com/xmdhs/hidethread/web"
)

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/", web.WebRoot)
	r.HandleFunc("/style.css", web.Style)
	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 12 * time.Second,
		Handler:      r,
	}
	log.Println(s.ListenAndServe())
}
