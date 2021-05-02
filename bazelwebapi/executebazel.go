package main

import (
	"os"
	"os/exec"
)

func (s *server) executeBazel(task *task) error {
	err := s.prepareGitRepo(task)
	if err != nil {
		return err
	}

	cmd := exec.Command("bazelisk", task.BazelCmd, task.Target)
	cmd.Dir = s.gitPath(task)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
