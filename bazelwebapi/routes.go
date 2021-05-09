package main

import "net/http"

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/execute/task", s.checkApiKeyExec(s.ExecuteTask))
	mux.Handle("/getstate/task/", s.checkApiKeyRead(s.GetStateOfTask))
	mux.Handle("/get/task/", s.checkApiKeyRead(s.GetTask))
	mux.HandleFunc("/health/ready", s.Ready)
	mux.HandleFunc("/health/live", s.Live)
	return mux
}
