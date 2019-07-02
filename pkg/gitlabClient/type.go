package gitlabClient

import (
	"github.com/xanzy/go-gitlab"
)

type GitLab struct {
	Client *gitlab.Client
	Config Config
}

type Config struct {
	Labels    []gitlab.Label
	BoardList []gitlab.BoardList
	Current
}

type Current struct {
	Sprint string
}

type Weight struct {
	Actually int
	Planned  int
}
