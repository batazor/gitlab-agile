package gitlabClient

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

func init() {
	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	viper.SetDefault("GITLAB_URL", "https://gitlab.com/api/v4")
	viper.SetDefault("GITLAB_TOKEN", "XXXXXXXXXXXXXXX")

	// Set config from .ENV
	if _, err := os.Stat(".env"); err == nil {
		config, err := ioutil.ReadFile(".env")
		if err != nil {
			logger.Error("Error read .env file", zap.Error(err))
		}

		viper.SetConfigType("env")
		viper.ReadConfig(bytes.NewBuffer(config))
	}
}

var GITLAB GitLab

func Run() (GitLab, error) {
	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	viper.AutomaticEnv()

	GITLAB.Client = gitlab.NewClient(nil, viper.GetString("GITLAB_TOKEN"))
	err := GITLAB.Client.SetBaseURL(viper.GetString("GITLAB_URL"))

	fmt.Println(viper.GetString("GITLAB_URL"))
	fmt.Println(viper.GetString("GITLAB_TOKEN"))

	return GITLAB, err
}
