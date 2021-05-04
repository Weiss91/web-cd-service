package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type server struct {
	sync.Mutex
	c           *config
	statusMap   map[string]*status
	runningTask string //uuid of running task
}

type task struct {
	Remote           string
	Commit           string
	Target           string
	BazelCmd         string // can be e.g. build, test, run
	Prio             int    // for later use priorisation of builds
	Registry         string
	RegistryUser     string
	RegistryPassword string
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
	id := uuid.NewString()
	now := time.Now()
	s.statusMap[id] = &status{
		id:      id,
		start:   now,
		updated: now,
		state:   WAITING,
	}

	for s.runningTask != "" {
		log.Println("waiting for execution of task ", id)
		time.Sleep(time.Second * 5)
	}

	s.runningTask = id
	// release running task in any case to not deadlock server
	defer s.releaseTask()
	s.statusMap[id].updated = time.Now()
	s.statusMap[id].state = RUNNING

	result, err := s.executeBazel(t)
	if err != nil {
		msg := fmt.Sprintf("error: %s\n\nbazel output:\n%s", err.Error(), result.String())
		writeError(w, http.StatusInternalServerError, msg)
		return
	}
	w.Write(result.Bytes())
}

func (s *server) releaseTask() {
	if s.statusMap != nil && s.statusMap[s.runningTask] != nil {
		log.Printf("Task %s done", s.runningTask)
		s.statusMap[s.runningTask].state = DONE
		now := time.Now()
		s.statusMap[s.runningTask].end = now
		s.statusMap[s.runningTask].updated = now
	}
	s.runningTask = ""
}
