package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type task struct {
	Remote       string
	Commit       string
	Target       string
	BazelCmd     string // can be e.g. build, test, run
	Prio         string // for priorisation of tasks
	Registry     string
	Output       string
	Id           string
	Start        time.Time
	End          time.Time
	Updated      time.Time
	State        STATE
	StateString  string
	Err          string
	PrioCategory PRIO
	PrioString   string
	// triggeredBy string/enum later when gitlab webhook or other sources are implemented
}

type tasks struct {
	sync.Mutex
	tasks map[string]*task
}

func newTasks() *tasks {
	return &tasks{
		tasks: make(map[string]*task),
	}
}

func (ts *tasks) add(t *task) {
	ts.Lock()
	defer ts.Unlock()
	ts.tasks[t.Id] = t
}

func (ts *tasks) find(id string) task {
	ts.Lock()
	defer ts.Unlock()
	val, ok := ts.tasks[id]
	if !ok {
		return task{}
	}
	return *val
}

func (ts *tasks) delete(id string) {
	ts.Lock()
	defer ts.Unlock()
	delete(ts.tasks, id)
}

func parseTask(r *http.Request) (*task, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read request body")
	}
	t := &task{}
	err = json.Unmarshal(b, t)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	t.PrioCategory = newPrio(t.Prio)
	t.PrioString = t.PrioCategory.ToString()
	return t, nil
}

func getTask(s *server, path string) (task, error) {
	pathparts := strings.Split(path, "/")
	uuid := pathparts[len(pathparts)-1]

	t := s.activeTasks.find(uuid)
	if t.Id == "" {
		t = s.history.find(uuid)
	}

	if t.StateString == "" {
		return t, fmt.Errorf("no task found with id %s", uuid)
	}
	return t, nil
}

func (t *task) setState(s STATE) {
	t.State = s
	t.StateString = s.ToString()
}
