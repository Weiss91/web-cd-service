package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"
)

func (s *server) executeBazel(task *task) (*bytes.Buffer, error) {
	err := s.prepareGitRepo(task)
	if err != nil {
		return nil, err
	}

	// create docker conf file before needed
	dc := &dockerconf{}
	if task.BazelCmd == "run" {
		if task.Registry != "" {
			dc.DockerConfPath = s.c.DockerConfPath
			dc.Registry = task.Registry
			dc.Auth = base64.StdEncoding.EncodeToString(
				[]byte(fmt.Sprintf("%s:%s", task.RegistryUser, task.RegistryPassword)))
		}
		err = dc.createDockerConf()
		if err != nil {
			return nil, err
		}
		// remove it in any case after executing this function
		defer dc.removeDockerConf()
	}

	cmd := exec.Command("bazelisk", task.BazelCmd, task.Target)
	cmd.Dir = s.gitPath(task)

	bufStdErr := bytes.Buffer{}
	bufStdOut := bytes.Buffer{}
	cmd.Stderr = &bufStdErr
	cmd.Stdout = &bufStdOut

	err = cmd.Run()
	if err != nil {
		return &bufStdErr, err
	}

	// append outputs
	bufStdErr.Write([]byte("\nStdout:\n"))
	bufStdErr.Write(bufStdOut.Bytes())
	return &bufStdErr, nil
}
