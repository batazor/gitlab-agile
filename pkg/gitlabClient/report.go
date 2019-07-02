package gitlabClient

import "fmt"

func (git *GitLab) ReportPlannedActually() error {
	weightTotal := Weight{
		Actually: 0,
		Planned:  0,
	}

	issues, err := git.GetMilestoneIssues(git.Config.Current.Sprint)
	if err != nil {
		return err
	}

	for _, issue := range issues {
		weight := git.ParseWeight(issue.Title)

		if issue.State == "closed" {
			weightTotal.Actually += weight
		}

		weightTotal.Planned += weight
	}

	fmt.Println("weightTotal", weightTotal)

	return nil
}
