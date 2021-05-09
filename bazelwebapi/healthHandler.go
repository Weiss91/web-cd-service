package main

import "net/http"

func (s *server) Ready(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeStatusNotAllowed(w)
		return
	}
	w.Write([]byte("Ready"))
}

func (s *server) Live(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeStatusNotAllowed(w)
		return
	}
	w.Write([]byte("Live"))
}
