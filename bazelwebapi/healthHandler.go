package main

import (
	"fmt"
	"net/http"
	"time"
)

func (s *server) Ready(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeStatusNotAllowed(w)
		return
	}
	w.Write([]byte(fmt.Sprintf("Ready since %.2f hours", time.Since(s.start).Hours())))
}

func (s *server) Live(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeStatusNotAllowed(w)
		return
	}
	w.Write([]byte(fmt.Sprintf("Live since %.2f hours", time.Since(s.start).Hours())))
}
