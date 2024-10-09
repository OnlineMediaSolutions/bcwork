package messager

import (
	"fmt"
	"log"

	"github.com/m6yf/bcwork/config"
	"github.com/slack-go/slack"
)

type Messager interface {
	SendMessage(message string) error
}

var _ Messager = (*SlackModule)(nil)

type SlackModule struct {
	api       *slack.Client
	channelID string
}

func NewSlackModule() (*SlackModule, error) {
	credentials, err := config.FetchConfigValues([]string{"slack_token", "slack_alerts_channel"})
	if err != nil {
		return nil, err
	}
	slackToken := credentials["slack_token"]
	channelID := credentials["slack_alerts_channel"]
	api := slack.New(slackToken)
	return &SlackModule{
		api:       api,
		channelID: channelID,
	}, nil
}

func (sm *SlackModule) SendMessage(message string) error {
	_, _, err := sm.api.PostMessage(
		sm.channelID,
		slack.MsgOptionText(message, false),
	)

	if err != nil {
		log.Printf("Error sending message to Slack: %v", err)
		return err
	}

	fmt.Println("Message sent successfully")
	return nil
}
