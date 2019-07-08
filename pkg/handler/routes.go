package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/nlopes/slack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func init() {
	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	viper.SetDefault("SLACK_TOKEN", "XXXXXXXXXXXXXXX")

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

// Routes creates a REST router
func Routes() chi.Router {
	r := chi.NewRouter()

	client := slack.New(viper.GetString("SLACK_TOKEN"))
	slackListener := &SlackListener{
		client: client,
	}
	go slackListener.ListenAndResponse()

	r.Post("/", ServeHTTP)

	return r
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed to unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Only accept message from slack with valid token
	//if message.Token != h.verificationToken {
	//	log.Printf("[ERROR] Invalid token: %s", message.Token)
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	action := message.ActionCallback.AttachmentActions[0]

	fmt.Println("action.Name", action.Name)

	switch action.Name {
	case actionSelect:
		value := action.SelectedOptions[0].Value

		switch strings.Title(value) {
		case "ReportByMilestoune":
			// Overwrite original drop down message.
			originalMessage := message.OriginalMessage
			originalMessage.Attachments[0].Text = "Enter command: `ReportByMilestoune [nameMilestoune]`"
			originalMessage.Attachments[0].Actions = []slack.AttachmentAction{}

			w.Header().Add("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&originalMessage)
			return
		}

		// Overwrite original drop down message.
		originalMessage := message.OriginalMessage
		originalMessage.Attachments[0].Text = fmt.Sprintf("OK to order %s ?", strings.Title(value))
		originalMessage.Attachments[0].Actions = []slack.AttachmentAction{
			{
				Name:  actionStart,
				Text:  "Yes",
				Type:  "button",
				Value: "start",
				Style: "primary",
			},
			{
				Name:  actionCancel,
				Text:  "No",
				Type:  "button",
				Style: "danger",
			},
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&originalMessage)
		return
	case actionStart:
		title := ":ok: your order was submitted! yay!"
		responseMessage(w, message.OriginalMessage, title, "")
		return
	case actionCancel:
		title := fmt.Sprintf(":x: @%s canceled the request", message.User.Name)
		responseMessage(w, message.OriginalMessage, title, "")
		return
	default:
		log.Printf("[ERROR] ]Invalid action was submitted: %s", message.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// responseMessage response to the original slackbutton enabled message.
// It removes button and replace it with message which indicate how bot will work
func responseMessage(w http.ResponseWriter, original slack.Message, title, value string) {
	//original.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	original.Attachments = append(original.Attachments, slack.Attachment{
		Title: title,
		Text:  value,
	})
	fmt.Println("original.Attachments", original.Attachments)

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&original)
}
