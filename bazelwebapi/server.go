package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

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
	ts.tasks[t.id] = t
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

type queue struct {
	sync.Mutex
	tasks []*task
}

func newQueue() *queue {
	return &queue{
		tasks: []*task{},
	}
}

func (q *queue) next() (t *task) {
	q.Lock()
	defer q.Unlock()
	if len(q.tasks) > 0 {
		t = q.tasks[0]
		q.tasks = q.tasks[1:]
	}
	return t
}

func (q *queue) add(t *task) {
	q.Lock()
	defer q.Unlock()

	q.tasks = append(q.tasks, t)

	// sort for 1. prio asc and 2. start time asc
	sort.Slice(q.tasks, func(i, j int) bool {
		if q.tasks[i].Prio < q.tasks[j].Prio {
			return true
		}
		if q.tasks[i].Prio > q.tasks[j].Prio {
			return false
		}
		return q.tasks[i].start.Before(q.tasks[j].start)
	})
}

type server struct {
	sync.Mutex
	c           *config
	history     *tasks
	activeTasks *tasks
	queue       *queue
	runningTask string //uuid of running task
}

type task struct {
	Remote   string
	Commit   string
	Target   string
	BazelCmd string // can be e.g. build, test, run
	Prio     string // for priorisation of tasks
	Registry string
	output   string
	id       string
	start    time.Time
	end      time.Time
	updated  time.Time
	state    STATE
	err      string
	prio     PRIO
	// triggeredBy string/enum later when gitlab webhook or other sources are implemented
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

	t.prio = newPrio(t.Prio)
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

	id := uuid.NewString()
	now := time.Now()
	t.id = id
	t.start = now
	t.updated = now
	t.state = WAITING

	s.activeTasks.add(t)
	s.queue.add(t)

	w.Write([]byte(fmt.Sprintf("{\"taskID\": \"%s\"}", id)))
}

func (s *server) GetStateOfTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeStatusNotAllowed(w)
		return
	}
	pathparts := strings.Split(r.URL.Path, "/")
	uuid := pathparts[len(pathparts)-1]

	t := s.activeTasks.find(uuid)
	if t.id == "" {
		t = s.history.find(uuid)
	}

	w.Write([]byte(fmt.Sprintf("{\"state\": \"%s\"}", t.state.ToString())))
}
