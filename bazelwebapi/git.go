package main

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func (s *server) getAuth() *http.BasicAuth {
	return &http.BasicAuth{
		Username: s.c.gitConf.User,
		Password: s.c.gitConf.Password,
	}
}

func (s *server) gitPath(task *task) string {
	pathparts := strings.Split(task.Remote, "/")
	return fmt.Sprintf("%s/%s", s.c.gitConf.Path, pathparts[len(pathparts)-1])
}

func (s *server) prepareGitRepo(task *task) error {
	r, err := git.PlainClone(s.gitPath(task), false,
		&git.CloneOptions{
			URL:  task.Remote,
			Auth: s.getAuth(),
		},
	)

	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			r, err = git.PlainOpen(s.gitPath(task))
			if err != nil {
				return err
			}
		}
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{Auth: s.getAuth(), Force: true})
	if err != nil {
		if err != git.NoErrAlreadyUpToDate && err != git.ErrNonFastForwardUpdate {
			return err
		}
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash:  plumbing.NewHash(task.Commit),
		Force: true,
	})
	if err != nil {
		return err
	}
	return nil
}
