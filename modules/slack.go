package modules

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

type SlackModule struct {
	api       *slack.Client
	channelID string
}

func NewSlackModule() *SlackModule {
	slackToken := viper.GetString("slack.token")
	channelID := viper.GetString("slack.channelID")
	api := slack.New(slackToken)
	return &SlackModule{
		api:       api,
		channelID: channelID,
	}
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

func (sm *SlackModule) SendAttachment(attachment slack.Attachment) error {
	_, _, err := sm.api.PostMessage(
		sm.channelID,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		log.Printf("Error sending attachement to Slack: %v", err)
		return err
	}

	fmt.Println("Message sent successfully")
	return nil
}
