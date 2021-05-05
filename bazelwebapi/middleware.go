package main

import (
	"net/http"
)

func (s *server) checkApiKeyExec(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isValidKey(w, r, s.c.apiKeyExec) {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *server) checkApiKeyRead(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isValidKey(w, r, s.c.apiKeyRead) {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isValidKey(w http.ResponseWriter, r *http.Request, key string) bool {
	k := r.Header.Get("X-API-KEY")
	if k != key {
		writeError(w, http.StatusUnauthorized, "invalid apikey")
		return false
	}
	return true
}
