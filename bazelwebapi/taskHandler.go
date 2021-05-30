package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *server) ExecuteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeStatusNotAllowed(w)
		return
	}

	t, err := parseTask(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id := uuid.NewString()
	now := time.Now()
	t.Id = id
	t.Start = now
	t.Updated = now
	t.setState(WAITING)

	s.activeTasks.add(t)
	s.saveActiveTasks()
	s.queue.add(t)

	w.Write([]byte(fmt.Sprintf("{\"taskID\": \"%s\"}", id)))
}

func (s *server) GetStateOfTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeStatusNotAllowed(w)
		return
	}
	t, err := getTask(s, r.URL.Path)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Write([]byte(fmt.Sprintf("{\"state\": \"%s\"}", t.State.ToString())))
}

func (s *server) GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeStatusNotAllowed(w)
		return
	}

	t, err := getTask(s, r.URL.Path)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	b, err := json.Marshal(t)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}

	w.Write(b)
}
