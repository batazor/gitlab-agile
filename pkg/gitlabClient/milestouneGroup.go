package gitlabClient

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

func (git *GitLab) GetMilestoneByGroupName(pid gitlab.Group, milestoune string) (*gitlab.GroupMilestone, error) {
	opt := gitlab.ListGroupMilestonesOptions{
		Search: milestoune,
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		// Get the first page with projects.
		milesounes, resp, err := git.Client.GroupMilestones.ListGroupMilestones(pid.ID, &opt)
		if err != nil {
			return nil, err
		}

		// List all the projects we've found so far.
		//for _, p := range milesounes {
		//	fmt.Printf("Found milesoune: %s\n", p.Title)
		//}

		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			if len(milesounes) != 0 {
				return milesounes[0], err
			}

			return nil, err
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}
}

func (git *GitLab) GetGroupMilestoneIssues(name string) ([]*gitlab.Issue, error) {
	var issues []*gitlab.Issue
	nameMilestoune := fmt.Sprintf("\"%s\"", name)

	projects, err := git.ListGroup()
	if err != nil {
		return nil, err
	}

	for _, group := range projects {
		m, err := git.GetMilestoneByGroupName(*group, nameMilestoune)

		if err != nil || m == nil {
			// TODO: Check status code 403/404/etc
			//return nil, err
			continue
		}

		opt := gitlab.GetGroupMilestoneIssuesOptions{
			PerPage: 100,
			Page:    1,
		}

		for {
			// Get the first page with projects.
			iss, resp, err := git.Client.GroupMilestones.GetGroupMilestoneIssues(group.ID, m.ID, &opt)
			if err != nil {
				return issues, err
			}

			// Exit the loop when we've seen all pages.
			if resp.CurrentPage >= resp.TotalPages {
				issues = append(issues, iss...)
				break
			}

			// Update the page number to get the next page.
			opt.Page = resp.NextPage
		}
	}

	return issues, nil
}
