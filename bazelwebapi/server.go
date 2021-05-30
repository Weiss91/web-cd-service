package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type server struct {
	sync.Mutex
	start       time.Time
	c           *config
	history     *tasks
	activeTasks *tasks
	queue       *queue
	runningTask string //uuid of running task
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func writeStatusNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
}

func (s *server) executeBazel(task *task) error {
	err := s.prepareGitRepo(task)
	if err != nil {
		return err
	}

	// create only registry when is set in task
	if task.Registry != "" {
		err = s.c.dockerConf.createDockerConf(task.Registry)
		if err != nil {
			return err
		}
		// remove it in any case after executing this function
		defer s.c.dockerConf.removeDockerConf()
	}

	cmd := exec.Command("bazelisk", task.BazelCmd, task.Target)
	cmd.Dir = s.gitPath(task)

	bufStdErr := bytes.Buffer{}
	bufStdOut := bytes.Buffer{}
	cmd.Stderr = &bufStdErr
	cmd.Stdout = &bufStdOut

	err = cmd.Run()
	if err != nil {
		task.Err = err.Error()
		task.Output = bufStdErr.String()
		return nil
	}

	// append outputs
	bufStdErr.Write([]byte("\nStdout:\n"))
	bufStdErr.Write(bufStdOut.Bytes())
	task.Output = bufStdErr.String()
	return nil
}

func (s *server) saveActiveTasks() error {
	path := filepath.Join(s.c.storageConf.Path, "active")
	s.activeTasks.Lock()
	defer s.activeTasks.Unlock()
	err := saveTasks(path, s.activeTasks)
	if err != nil {
		return err
	}
	log.Println("Successfull saved active tasks")
	return nil
}

func (s *server) loadActiveTasks() error {
	path := filepath.Join(s.c.storageConf.Path, "active")
	ts, err := loadTasks(path)
	if err != nil {
		return err
	}

	s.activeTasks.Lock()
	defer s.activeTasks.Unlock()
	s.activeTasks = ts

	log.Println("Successfull loaded active tasks")
	return nil
}

func shutdownBazelServer() {
	cmd := exec.Command("bazelisk", "shutdown")
	cmd.Run()
	cmd.Stdout = os.Stdout
}

func (s *server) prepareShutdown() {
	go shutdownBazelServer()
	s.saveActiveTasks()
	time.Sleep(time.Second * 10)
}
