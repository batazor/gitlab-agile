package handler

import (
	"encoding/csv"
	"fmt"
	"github.com/batazor/gitlab-agile/pkg/gitlabClient"
	"github.com/nlopes/slack"
	"log"
	"os"
	"strconv"
	"strings"
)

type SlackListener struct {
	client *slack.Client
}

// LstenAndResponse listens slack events and response
// particular messages. It replies by slack message button.
func (s *SlackListener) ListenAndResponse() {
	rtm := s.client.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

const (
	// action is used for slack attament action.
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)

// handleMesageEvent handles message events.
func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[0:]
	fmt.Println("ev.Channel", m)
	//if len(m) == 0 || m[0] != "hey" {
	//	return fmt.Errorf("invalid message")
	//}

	attachment := slack.Attachment{}

	switch m[0] {
	case "ReportByMilestoune":
		nameMilestone := strings.Join(m[1:], " ")

		weightTotal, err := gitlabClient.GITLAB.ReportPlannedActually(nameMilestone)
		if err != nil {
			fmt.Println("Error", err)
		}

		var data = [][]string{{"Planned", strconv.Itoa(weightTotal.Planned)}, {"Actually", strconv.Itoa(weightTotal.Actually)}}
		file, err := os.Create("result.csv")
		if err != nil {
			fmt.Println("Error", err)
		}

		writer := csv.NewWriter(file)

		for _, value := range data {
			err := writer.Write(value)
			if err != nil {
				fmt.Println("Error", err)
			}
		}

		writer.Flush()
		file.Close()

		params := slack.FileUploadParameters{
			Filename: file.Name(),
			Title:    file.Name(),
			Channels: []string{ev.Channel},
			File:     file.Name(),
		}

		_, err = s.client.UploadFile(params)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		return nil
	case "ReportByMilestouneUser":
		nameMilestone := strings.Join(m[1:], " ")

		weightTotals, err := gitlabClient.GITLAB.ReportPlannedActuallyByUser(nameMilestone)
		if err != nil {
			fmt.Println("Error", err)
		}

		var data = [][]string{}

		header := []string{"Username", "Planned", "Actually"}
		data = append(data, header)

		for _, weightTotal := range weightTotals {
			header := []string{weightTotal.Nickname, strconv.Itoa(weightTotal.Planned), strconv.Itoa(weightTotal.Actually)}
			data = append(data, header)
		}

		file, err := os.Create("result.csv")
		if err != nil {
			fmt.Println("Error", err)
		}

		writer := csv.NewWriter(file)

		for _, value := range data {
			err := writer.Write(value)
			if err != nil {
				fmt.Println("Error", err)
			}
		}

		writer.Flush()
		file.Close()

		params := slack.FileUploadParameters{
			Filename: file.Name(),
			Title:    file.Name(),
			Channels: []string{ev.Channel},
			File:     file.Name(),
		}

		_, err = s.client.UploadFile(params)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		return nil
	case "reportByIssue":
		nameMilestone := strings.Join(m[1:], " ")

		var data = [][]string{}

		header := []string{
			"Title",
			"Weight",
			"Milestone",
			"Milestone start",
			"Milestone end",
			"Author",
			"State",
			"Assignee",
			"Labels",
			"CreatedAt",
			"UpdatedAt",
			"ClosedAt",
			"DueDate",
			"WebURL",
		}
		data = append(data, header)

		issues, err := gitlabClient.GITLAB.ListIssue(&nameMilestone)
		if err != nil {
			fmt.Println("Error", err)
		}

		for i, issue := range issues {
			weight := gitlabClient.GITLAB.ParseWeight(issue.Title)

			var DueDate string
			if issues[i].DueDate != nil {
				DueDate = issues[i].DueDate.String()
			}

			header := []string{
				issue.Title,
				strconv.Itoa(weight),
				issue.Milestone.Title,
				issue.Milestone.StartDate.String(),
				issue.Milestone.DueDate.String(),
				issue.Author.Username,
				issue.State,
				issue.Assignee.Username,
				strings.Join(issue.Labels, ","),
				issues[i].CreatedAt.String(),
				issues[i].UpdatedAt.String(),
				issues[i].ClosedAt.String(),
				DueDate,
				issue.WebURL,
			}
			data = append(data, header)
		}

		file, err := os.Create("result.csv")
		if err != nil {
			fmt.Println("Error", err)
		}

		writer := csv.NewWriter(file)

		for _, value := range data {
			err := writer.Write(value)
			if err != nil {
				fmt.Println("Error", err)
			}
		}

		writer.Flush()
		file.Close()

		params := slack.FileUploadParameters{
			Filename: file.Name(),
			Title:    file.Name(),
			Channels: []string{ev.Channel},
			File:     file.Name(),
		}

		_, err = s.client.UploadFile(params)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		return nil
	case "hey":
		// value is passed to message handler when request is approved.
		attachment = slack.Attachment{
			Text:       "Select option?",
			Color:      "#f9a41b",
			CallbackID: "selectOption",
			Actions: []slack.AttachmentAction{
				{
					Name: actionSelect,
					Type: "select",
					Options: []slack.AttachmentActionOption{
						{
							Text:  "Milestoune",
							Value: "Milestoune",
						},
						{
							Text:  "Report by milestoune",
							Value: "reportByMilestoune",
						},
						{
							Text:  "Report by milestoune && user",
							Value: "reportByMilestouneUser",
						},
						{
							Text:  "Report by issues",
							Value: "reportByIssue",
						},
					},
				},

				{
					Name:  actionCancel,
					Text:  "Cancel",
					Type:  "button",
					Style: "danger",
				},
			},
		}
	default:
		return nil
	}

	if _, _, err := s.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment)); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}
