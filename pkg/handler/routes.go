package handler

import (
	"bytes"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
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
	//r.Use(middleware.AllowContentType("application/json"))

	r.Post("/", Example)

	return r
}

func Example(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//api := slack.New(viper.GetString("SLACK_TOKEN"), slack.OptionDebug(true))

	str := `{
            "text": "This is your first interactive message",
            "attachments": [
                {
                    "text": "Building buttons is easy right?",
                    "fallback": "Shame... buttons aren't supported in this land",
                    "callback_id": "button_tutorial",
                    "color": "#3AA3E3",
                    "attachment_type": "default",
                    "actions": [
                        {
                            "name": "yes",
                            "text": "yes",
                            "type": "button",
                            "value": "yes"
                        },
                        {
                            "name": "no",
                            "text": "no",
                            "type": "button",
                            "value": "no"
                        },
                        {
                            "name": "maybe",
                            "text": "maybe",
                            "type": "button",
                            "value": "maybe",
                            "style": "danger"
                        }
                    ]
                }
            ]
        }`

	// Normal JSON stuff
	w.Write([]byte(str))
}
