package main

import "strings"

type PRIO int

const (
	HOTFIX PRIO = iota + 1
	PROD
	INTEGRATIONANDTESTING
	DEVELOP
	UNIMPORTANT
)

func (i PRIO) ToString() string {
	switch i {
	case HOTFIX:
		return "hotfix"
	case PROD:
		return "prod"
	case INTEGRATIONANDTESTING:
		return "iat"
	case DEVELOP:
		return "develop"
	case UNIMPORTANT:
		return "unimportant"
	default:
		return "not implemented"
	}
}

func newPrio(s string) PRIO {
	switch strings.ToLower(s) {
	case "hotfix", "bugfix", "1":
		return HOTFIX
	case "prod", "production", "2":
		return PROD
	case "integrationandtesting", "iat", "testing", "test", "3":
		return INTEGRATIONANDTESTING
	case "dev", "develop", "feature", "4":
		return DEVELOP
	default:
		return UNIMPORTANT
	}
}
