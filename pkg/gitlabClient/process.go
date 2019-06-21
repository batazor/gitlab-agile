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
	process := Process{}
	err = yaml.Unmarshal([]byte(config), &process)
	if err != nil {
		return err
	}

	// CONFIG
	LABELS := process.Labels

	// Get list group
	groups, err := git.ListGroup()
	for _, group := range groups {
		// Create labels for each project
		for _, label := range LABELS {
			err := git.CreateGroupLabel(group.ID, label)
			if err != nil {
				logger.Info(
					"Create labels for group",
					zap.Int("group Id", group.ID),
					zap.String("label Name", label.Name),)
			}
		}
	}

	return nil
}
