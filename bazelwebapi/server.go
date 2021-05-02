package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type server struct {
	sync.Mutex
	c         *config
	statusMap map[string]status
}

type task struct {
	Remote   string
	Commit   string
	Target   string
	BazelCmd string // can be e.g. build, test, run
	Prio     int    // for later use priorisation of builds
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func writeStatusNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
}

func getTask(r *http.Request) (*task, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read request body")
	}
	t := &task{}
	err = json.Unmarshal(b, t)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return t, nil
}

func (s *server) ExecuteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeStatusNotAllowed(w)
		return
	}

	t, err := getTask(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Repo: %s;\nCommit: %s\nExecuting Target: %s\nwith command: %s\n", t.Remote, t.Commit, t.Target, t.BazelCmd)

	s.Lock()
	defer s.Unlock()

	err = s.executeBazel(t)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
