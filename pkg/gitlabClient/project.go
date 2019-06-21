package gitlabClient

import (
	"github.com/xanzy/go-gitlab"
)

func (git *GitLab) ListProject() ([]*gitlab.Project, error) {
	FALSE := false
	TRUE := true

	opt := &gitlab.ListProjectsOptions{
		Archived: &FALSE,
		Membership: &TRUE,
	}
	projects, _, err := git.Client.Projects.ListProjects(opt)

	return projects, err
}

func (git *GitLab) CreateProjectLabel(projectId interface{}, label gitlab.Label) error {
	// Create new label
	l := &gitlab.CreateLabelOptions{
		Name:  &label.Name,
		Color: &label.Color,
	}
	_, _, err := git.Client.Labels.CreateLabel(projectId, l)
	if err != nil {
		return err
	}

	return nil
}