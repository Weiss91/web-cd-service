package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"text/template"
)

type config struct {
	GitUser        string
	GitPassword    string
	GitDataPath    string
	Registry       string
	ApiKeyExec     string
	ApiKeyRead     string
	DockerConfPath string
	Auth           string
}

type dockerconf struct {
	Registry       string
	Auth           string
	DockerConfPath string
}

func loadconfig() (*config, error) {
	gitUser := os.Getenv("GIT_USER")
	gitPw := os.Getenv("GIT_PASSWORD")
	gitDataPath := os.Getenv("GIT_DATA_PATH")
	registry := os.Getenv("REGISTRY")
	registryUser := os.Getenv("REGISTRY_USER")
	registryPw := os.Getenv("REGISTRY_PASSWORD")
	dockerConfPath := os.Getenv("DOCKER_CONF_PATH")
	apiKeyRead := os.Getenv("API_KEY_READ")
	apiKeyExec := os.Getenv("API_KEY_EXEC")

	if gitUser == "" {
		return nil, fmt.Errorf("GIT_USER not set")
	}

	if gitPw == "" {
		return nil, fmt.Errorf("GIT_PASSWORD not set")
	}

	if gitDataPath == "" {
		return nil, fmt.Errorf("GIT_DATA_PATH not set")
	}

	return &config{
		GitUser:        gitUser,
		GitPassword:    gitPw,
		GitDataPath:    gitDataPath,
		ApiKeyRead:     apiKeyRead,
		ApiKeyExec:     apiKeyExec,
		Registry:       registry,
		DockerConfPath: dockerConfPath,
		Auth:           base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", registryUser, registryPw))),
	}, nil
}

func (c *dockerconf) createDockerConf() error {
	tmpl, err := template.New("dockerconf").Parse(dockerconftmpl)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(c.DockerConfPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	err = tmpl.Execute(f, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *dockerconf) removeDockerConf() error {
	return os.Remove(c.DockerConfPath)
}

const dockerconftmpl = `
{
	"auths": {
		"{{.Registry}}": {
			"auth": "{{.Auth}}"
		}
	}
}
`
