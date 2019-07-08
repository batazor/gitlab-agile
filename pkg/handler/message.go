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
