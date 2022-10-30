package gateway

import "github.com/slack-go/slack"

func SlackGetUserInfo(token string, slack_user_id string) (*slack.User, error) {
	api := slack.New(token)
	user, err := api.GetUserInfo(slack_user_id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
