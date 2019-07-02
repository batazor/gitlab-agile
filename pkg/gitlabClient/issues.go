package gitlabClient

import (
	"fmt"
	"regexp"
	"strconv"
)

func (git *GitLab) ParseWeight(title string) int {
	var rgx = regexp.MustCompile(`\[(.*?)\]`)
	rs := rgx.FindStringSubmatch(title)

	fmt.Println(rs[1])

	if len(rs) == 0 {
		return 0
	}

	weight, err := strconv.Atoi(rs[1])
	if err != nil {
		return 0
	}

	return weight
}
