package gitlabClient

import (
	"github.com/xanzy/go-gitlab"
)

type GitLab struct {
	Client *gitlab.Client
}

type Process struct {
	Labels    []gitlab.Label
	BoardList []gitlab.BoardList
}

type Weight struct {
	Actually int
	Planned  int
}
