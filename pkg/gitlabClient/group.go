package gitlabClient

import (
	"github.com/xanzy/go-gitlab"
)

func (git *GitLab) ListGroup() ([]*gitlab.Group, error) {
	opt := &gitlab.ListGroupsOptions{}
	groups, _, err := git.Client.Groups.ListGroups(opt)

	return groups, err
}

func (git *GitLab) CreateGroupLabel(projectId interface{}, label gitlab.Label) error {
	// Create new label
	l := &gitlab.CreateGroupLabelOptions{
		Name:  &label.Name,
		Color: &label.Color,
	}
	_, _, err := git.Client.GroupLabels.CreateGroupLabel(projectId, l)
	if err != nil {
		return err
	}

	return nil
}
