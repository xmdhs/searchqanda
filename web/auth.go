package web

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

func Auth(HandleFunc func(http.ResponseWriter, *http.Request), password string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Query().Get("password")
		if p == "" {
			c, err := r.Cookie("password")
			if err != nil {
				w.WriteHeader(403)
				return
			}
			b, err := hex.DecodeString(c.Value)
			if err != nil {
				w.WriteHeader(403)
				return
			}
			p = string(b)
		}
		if tosha256(p) == password {
			hex := hex.EncodeToString([]byte(p))
			http.SetCookie(w, &http.Cookie{
				Name:     "password",
				Value:    hex,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})
			HandleFunc(w, r)
		} else {
			w.WriteHeader(403)
			return
		}
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