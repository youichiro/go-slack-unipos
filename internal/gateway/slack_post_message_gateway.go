package gateway

import (
	"os"

	"github.com/slack-go/slack"
)

func SlackPostMessage(token string, message string) error {
	api := slack.New(token)
	_, _, err := api.PostMessage(
		os.Getenv("SLACK_UNIPOS_CHANNEL_ID"),
		slack.MsgOptionText(message, false),
		slack.MsgOptionEnableLinkUnfurl(),
	)
	return err
}
