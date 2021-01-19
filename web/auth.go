package web

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

func Auth(HandleFunc func(http.ResponseWriter, *http.Request), password string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if ok {
			p := tosha256(password)
			if p == password {
				HandleFunc(w, r)
			}
		}
		w.WriteHeader(401)
		w.Header().Add("www-authenticate", `Basic realm="password"`)
	}
}

func tosha256(s string) string {
	hash := sha256.New()
	_, err := hash.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
