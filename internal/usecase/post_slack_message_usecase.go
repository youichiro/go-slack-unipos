package usecase

import (
	"fmt"

	"github.com/youichiro/go-slack-my-unipos/internal/gateway"
)

func PostSlackMessageUsecase(slackToken string, senderSlackUserId string, slackUserIDs []string, message string, point int) error {
	mentionMsg := ""
	for _, slackUserID := range slackUserIDs {
		mentionMsg += "<@" + slackUserID + ">"
	}
	msg := fmt.Sprintf("from: <@%s>, to: %s, point: %d, message: %s", senderSlackUserId, mentionMsg, point, message)

	// ユーザーアイコンを取得する
	senderUser, err := gateway.SlackGetUserInfo(slackToken, senderSlackUserId)
	if err != nil {
		return err
	}

	err = gateway.SlackPostMessage(slackToken, msg, senderUser)
	if err != nil {
		return err
	}
	return nil
}
