package gitlabClient

func (git *GitLab) ReportPlannedActually(name string) (Weight, error) {
	weightTotal := Weight{
		Actually: 0,
		Planned:  0,
	}

	issues, err := git.GetMilestoneIssues(name)
	if err != nil {
		return weightTotal, err
	}

	for _, issue := range issues {
		weight := git.ParseWeight(issue.Title)

		if issue.State == "closed" {
			weightTotal.Actually += weight
		}

		weightTotal.Planned += weight
	}

	return weightTotal, nil
}
