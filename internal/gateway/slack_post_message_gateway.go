package gateway

import (
	"os"

	"github.com/slack-go/slack"
)

func SlackPostMessage(token string, message string, user *slack.User) error {
	api := slack.New(token)
	_, _, err := api.PostMessage(
		os.Getenv("SLACK_UNIPOS_CHANNEL_ID"),
		slack.MsgOptionText(message, false),
		slack.MsgOptionEnableLinkUnfurl(),
		slack.MsgOptionIconURL(user.Profile.Image48),
		slack.MsgOptionUsername(user.Profile.RealName),
	)
	return err
}
