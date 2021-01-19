package web

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

func Auth(HandleFunc func(http.ResponseWriter, *http.Request), password string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, p, ok := r.BasicAuth()
		if ok {
			p := tosha256(p)
			if p == password {
				http.SetCookie(w, &http.Cookie{
					Name:     "password",
					Value:    hex.EncodeToString([]byte(p)),
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				})
				HandleFunc(w, r)
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="password"`)
		w.WriteHeader(401)
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
