package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"
)

type config struct {
	gitConf     *gitConf
	dockerConf  *dockerConf
	storageConf *storageConf
	apiKeyExec  string
	apiKeyRead  string
	serverPort  string
}

type storageConf struct {
	Path string
}

type gitConf struct {
	User     string
	Password string
	Path     string
}

type dockerConf struct {
	Path       string               `json:"path"`
	Registries map[string]*registry `json:"registries"`
}

type registry struct {
	Registry string
	User     string `json:"user"`
	Password string `json:"password"`
	Auth     string `json:"auth"` //generated when not set
}

func loadconfig() (*config, error) {
	gitUser := os.Getenv("GIT_USER")
	gitPw := os.Getenv("GIT_PASSWORD")
	gitPath := os.Getenv("GIT_DATA_PATH") // defaults to /git
	// if dockerConf.Path not set it will default to /git/docker
	docker := os.Getenv("DOCKER_CONF")
	// defaults to /storage
	storagePath := os.Getenv("STORAGE_PATH")

	apiKeyRead := os.Getenv("API_KEY_READ")
	apiKeyExec := os.Getenv("API_KEY_EXEC")
	serverPort := os.Getenv("SERVER_PORT")

	if serverPort == "" {
		serverPort = "8088"
	}

	if gitPath == "" {
		gitPath = "/git"
		log.Println("default git path to /git")
	}

	if gitUser == "" {
		return nil, fmt.Errorf("GIT_USER not set")
	}

	if gitPw == "" {
		return nil, fmt.Errorf("GIT_PASSWORD not set")
	}

	if docker == "" {
		log.Println("WARNING: no registry information set. Images will not be pushed")
	}

	if storagePath == "" {
		log.Println("default storage path to /storage")
		storagePath = "/storage"
	}

	dc := &dockerConf{}
	if docker != "" {
		err := json.Unmarshal([]byte(docker), dc)
		if err != nil {
			return nil, err
		}

		if dc.Path == "" {
			dc.Path = "/git/docker"
		}

		for k, v := range dc.Registries {
			if v.Auth == "" {
				v.Auth = base64.StdEncoding.EncodeToString(
					[]byte(fmt.Sprintf("%s:%s", v.User, v.Password)))
			}
			if v.Registry == "" {
				v.Registry = k
			}
		}
	}

	sc := &storageConf{
		Path: storagePath,
	}

	return &config{
		gitConf: &gitConf{
			User:     gitUser,
			Password: gitPw,
			Path:     gitPath,
		},
		serverPort:  serverPort,
		apiKeyRead:  apiKeyRead,
		apiKeyExec:  apiKeyExec,
		dockerConf:  dc,
		storageConf: sc,
	}, nil
}

func (c *dockerConf) createDockerConf(registry string) error {
	if len(c.Registries) == 0 {
		return fmt.Errorf("no registry information set on server")
	}
	if _, ok := c.Registries[registry]; !ok {
		return fmt.Errorf("no registry information found on server for registry %s", registry)
	}
	tmpl, err := template.New("dockerconf").Parse(dockerconftmpl)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(c.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	err = tmpl.Execute(f, c.Registries[registry])
	if err != nil {
		return err
	}
	return nil
}

func (c *dockerConf) removeDockerConf() error {
	return os.Remove(c.Path)
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
