package main

import "net/http"

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/execute/task", s.checkApiKeyExec(s.ExecuteTask))
	mux.Handle("/getstate/task/", s.checkApiKeyRead(s.GetStateOfTask))

	return mux
}
