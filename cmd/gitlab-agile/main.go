package main

import (
	"github.com/batazor/gitlab-agile/pkg/gitlabClient"
	"go.uber.org/zap"
)

func main() {
	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Run GitLab
	git, err := gitlabClient.Run()
	if err != nil {
		zap.Error(err)
		return
	}

	logger.Info("Run GitLab")

	// Create SCRUM-structure
	err = git.Apply()
	if err != nil {
		zap.Error(err)
	}

	// Create report
	err = git.ReportPlannedActually()
	if err != nil {
		zap.Error(err)
		return
	}
}
