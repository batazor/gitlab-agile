package gitlabClient

import (
	"github.com/xanzy/go-gitlab"
)

func (git *GitLab) ReportPlannedActually(name string) (Weight, error) {
	issues := []*gitlab.Issue{}
	weightTotal := Weight{
		Actually: 0,
		Planned:  0,
	}

	// from project
	projectIssues, err := git.GetProjectMilestoneIssues(name)
	if err != nil {
		return weightTotal, err
	}
	issues = append(issues, projectIssues...)

	// from groups
	groupIssues, err := git.GetGroupMilestoneIssues(name)
	if err != nil {
		return weightTotal, err
	}
	issues = append(issues, groupIssues...)

	for _, issue := range issues {
		weight := git.ParseWeight(issue.Title)

		if issue.State == "closed" {
			weightTotal.Actually += weight
		}

		weightTotal.Planned += weight
	}

	return weightTotal, nil
}

func (git *GitLab) ReportPlannedActuallyByUser(name string) (map[string]*Weight, error) {
	issues := []*gitlab.Issue{}
	weightTotal := map[string]*Weight{}

	// get users
	users, err := git.ListUser()
	if err != nil {
		return weightTotal, err
	}

	for _, user := range users {
		weightTotal[user.Username] = &Weight{
			Nickname: user.Username,
		}
	}

	//from project
	projectIssues, err := git.GetProjectMilestoneIssues(name)
	if err != nil {
		return weightTotal, err
	}
	issues = append(issues, projectIssues...)

	// from groups
	groupIssues, err := git.GetGroupMilestoneIssues(name)
	if err != nil {
		return weightTotal, err
	}
	issues = append(issues, groupIssues...)

	for _, issue := range issues {
		weight := git.ParseWeight(issue.Title)

		if weightTotal[issue.Assignee.Username] != nil {
			weightTotal[issue.Assignee.Username].Planned += weight

			if issue.State == "closed" {
				weightTotal[issue.Assignee.Username].Actually += weight
			}
		}
	}

	return weightTotal, nil
}
