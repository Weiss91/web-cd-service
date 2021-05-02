package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

type config struct {
	GitUser        string
	GitPassword    string
	GitDataPath    string
	Registry       string
	DockerConfPath string
	Auth           string
}

func loadconfig() (*config, error) {
	gitUser := os.Getenv("GIT_USER")
	gitPw := os.Getenv("GIT_PASSWORD")
	gitDataPath := os.Getenv("GIT_DATA_PATH")
	registry := os.Getenv("REGISTRY")
	registryUser := os.Getenv("REGISTRY_USER")
	registryPw := os.Getenv("REGISTRY_PASSWORD")
	dockerConfPath := os.Getenv("DOCKER_CONF_PATH")

	if gitUser == "" {
		return nil, fmt.Errorf("GIT_USER not set")
	}

	if gitPw == "" {
		return nil, fmt.Errorf("GIT_PASSWORD not set")
	}

	if gitDataPath == "" {
		return nil, fmt.Errorf("GIT_DATA_PATH not set")
	}

	if registry == "" {
		return nil, fmt.Errorf("REGISTRY not set")
	}

	if registryUser == "" {
		return nil, fmt.Errorf("REGISTRY_USER not set")
	}

	if registryPw == "" {
		return nil, fmt.Errorf("REGISTRY_PASSWORD not set")
	}

	if dockerConfPath == "" {
		return nil, fmt.Errorf("DOCKER_CONF_PATH not set")
	}

	return &config{
		GitUser:        gitUser,
		GitPassword:    gitPw,
		GitDataPath:    gitDataPath,
		Registry:       registry,
		DockerConfPath: dockerConfPath,
		Auth:           base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", registryUser, registryPw))),
	}, nil
}
