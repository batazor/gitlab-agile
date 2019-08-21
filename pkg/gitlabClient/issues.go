package gitlabClient

import (
	"github.com/xanzy/go-gitlab"
	"regexp"
	"strconv"
)

func (git *GitLab) ParseWeight(title string) int {
	var rgx = regexp.MustCompile(`\[(.*?)\]`)
	rs := rgx.FindStringSubmatch(title)

	if len(rs) == 0 {
		return 0
	}

	weight, err := strconv.Atoi(rs[1])
	if err != nil {
		return 0
	}

	return weight
}

func (git *GitLab) ListIssue(milestone *string) ([]*gitlab.Issue, error) {
	opt := &gitlab.ListIssuesOptions{
		Milestone: milestone,
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		// Get the first page with projects.
		issues, resp, err := git.Client.Issues.ListIssues(opt)
		if err != nil {
			return issues, err
		}

		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			return issues, err
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}
}
