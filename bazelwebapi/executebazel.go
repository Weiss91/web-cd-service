package main

import (
	"bytes"
	"os/exec"
)

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
		task.err = err.Error()
		task.output = bufStdErr.String()
		return nil
	}

	// append outputs
	bufStdErr.Write([]byte("\nStdout:\n"))
	bufStdErr.Write(bufStdOut.Bytes())
	task.output = bufStdErr.String()
	return nil
}
