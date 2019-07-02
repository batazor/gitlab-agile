package gitlabClient

import (
	"github.com/xanzy/go-gitlab"
)

func (git *GitLab) CreateMilestone(name string) error {
	return nil
}

func (git *GitLab) GetMilestoneByName(pid gitlab.Project, milestoune string) (*gitlab.Milestone, error) {
	opt := gitlab.ListMilestonesOptions{
		Search: milestoune,
	}

	milesounes, _, err := git.Client.Milestones.ListMilestones(pid.ID, &opt)

	if len(milesounes) != 0 {
		return milesounes[0], err
	}

	return nil, err
}

func (git *GitLab) GetMilestoneIssues(name string) ([]*gitlab.Issue, error) {
	var issues []*gitlab.Issue

	projects, err := git.GetListProject()
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		m, err := git.GetMilestoneByName(*project, name)
		if err != nil || m == nil {
			// TODO: Check status code 403/404/etc
			//return nil, err
			continue
		}

		opt := gitlab.GetMilestoneIssuesOptions{}

		iss, _, err := git.Client.Milestones.GetMilestoneIssues(project.ID, m.ID, &opt)
		if err != nil {
			return nil, err
		}

		for _, is := range iss {
			issues = append(issues, is)
		}
	}

	return issues, nil
}
