package main

import (
	"bytes"
	"net/http"
	"os/exec"
	"sync"
)

type server struct {
	sync.Mutex
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
