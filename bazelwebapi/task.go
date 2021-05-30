package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	Tasks map[string]*task
}

func newTasks() *tasks {
	return &tasks{
		Tasks: make(map[string]*task),
	}
}

func (ts *tasks) add(t *task) {
	ts.Lock()
	defer ts.Unlock()
	ts.Tasks[t.Id] = t
}

func (ts *tasks) find(id string) task {
	ts.Lock()
	defer ts.Unlock()
	val, ok := ts.Tasks[id]
	if !ok {
		return task{}
	}
	return *val
}

func (ts *tasks) delete(id string) {
	ts.Lock()
	defer ts.Unlock()
	delete(ts.Tasks, id)
}

func appendTask(path string, task *task) error {
	ts, err := loadTasks(path)
	if err != nil {
		return err
	}
	ts.add(task)
	err = saveTasks(path, ts)
	if err != nil {
		return err
	}
	return nil
}

func saveTasks(path string, ts *tasks) error {
	path = path + ".json"
	tempPath := path + "_temp.json"
	b, err := json.Marshal(ts)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(tempPath, b, 0644)
	if err != nil {
		return err
	}
	err = copy(tempPath, path)
	if err != nil {
		return err
	}
	err = os.Remove(tempPath)
	if err != nil {
		return err
	}
	return nil
}

func loadTasks(path string) (*tasks, error) {
	path = path + ".json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println("No Tasks found in file:", path)
		return newTasks(), nil
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return newTasks(), err
	}

	ts := &tasks{}
	err = json.Unmarshal(b, ts)
	if err != nil {
		return newTasks(), err
	}
	return ts, nil
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
