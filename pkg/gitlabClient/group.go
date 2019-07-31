package gitlabClient

import (
	"github.com/xanzy/go-gitlab"
)

func (git *GitLab) ListGroup() ([]*gitlab.Group, error) {
	opt := &gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		// Get the first page with projects.
		groups, resp, err := git.Client.Groups.ListGroups(opt)
		if err != nil {
			return groups, err
		}

		// List all the projects we've found so far.
		//for _, g := range groups {
		//	fmt.Printf("Found groups: %s", g.Name)
		//}

		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			return groups, err
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}
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
