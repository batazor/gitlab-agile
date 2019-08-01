package gitlabClient

import "github.com/xanzy/go-gitlab"

func (git *GitLab) ListUser() ([]*gitlab.User, error) {
	FALSE := false
	TRUE := true

	opt := &gitlab.ListUsersOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		Active:  &TRUE,
		Blocked: &FALSE,
	}

	for {
		// Get the first page with projects.
		users, resp, err := git.Client.Users.ListUsers(opt)
		if err != nil {
			return users, err
		}

		// List all the projects we've found so far.
		//for _, g := range groups {
		//	fmt.Printf("Found groups: %s", g.Name)
		//}

		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			return users, err
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}
}
