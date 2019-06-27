package gitlabClient

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

func (git *GitLab) CreateBoardList(groupId interface{}, board []gitlab.BoardList) error {
	t := 11063479

	// Create new label
	l := &gitlab.CreateGroupIssueBoardListOptions{
		LabelID: &t,
		//Name:  &board.Label.Name,
	}

	fmt.Println("groupId", groupId)
	_, _, err := git.Client.GroupIssueBoards.CreateGroupIssueBoardList(groupId, 516008, l)
	fmt.Println("err", err)
	if err != nil {
		return err
	}

	return nil
}
