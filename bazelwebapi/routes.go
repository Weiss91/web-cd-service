package main

import "net/http"

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/execute/task", s.checkApiKeyExec(s.ExecuteTask))

	return mux
}
