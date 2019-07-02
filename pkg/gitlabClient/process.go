package gitlabClient

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Apply - read `process.yaml` and create:
// create: labels/file in repository/ETC
func (git *GitLab) Apply() error {
	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Read process config
	config, err := ioutil.ReadFile("process.yaml")
	if err != nil {
		return err
	}

	// Parse config
	err = yaml.Unmarshal([]byte(config), &git.Config)
	if err != nil {
		logger.Info(err.Error())
		return err
	}

	// CONFIG
	LABELS := git.Config.Labels
	BOARD := git.Config.BoardList

	// Get list group
	groups, err := git.ListGroup()
	for _, group := range groups {
		// Create labels for each group
		if false {
			for _, label := range LABELS {
				err := git.CreateGroupLabel(group.ID, label)
				if err != nil {
					logger.Info(
						"Create labels for group",
						zap.Int("group Id", group.ID),
						zap.String("label Name", label.Name))
				}
			}
		}

		// Create board for each group
		if false {
			err := git.CreateBoardList(group.ID, BOARD)
			if err != nil {
				logger.Info(
					"Create board list for group",
					zap.Int("group Id", group.ID),
				)
			}
		}
	}

	return nil
}
