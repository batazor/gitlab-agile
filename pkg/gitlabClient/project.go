package gitlabClient

import (
	"github.com/xanzy/go-gitlab"
)

func (git *GitLab) GetListProject() ([]*gitlab.Project, error) {
	FALSE := false
	TRUE := true

	opt := &gitlab.ListProjectsOptions{
		Archived:   &FALSE,
		Membership: &TRUE,
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		// Get the first page with projects.
		projects, resp, err := git.Client.Projects.ListProjects(opt)
		if err != nil {
			return projects, err
		}

		// List all the projects we've found so far.
		//for _, p := range projects {
		//	fmt.Printf("Found project: %s", p.Name)
		//}

		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			return projects, err
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}
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
