package usecase

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/youichiro/go-slack-my-unipos/internal/gateway"
	"github.com/youichiro/go-slack-my-unipos/internal/util"
)

func generateMessageBlocks(senderUser *slack.User, targetUsers []*slack.User, message string, point string) []slack.Block {
	dearBlockElements := make([]slack.MixedElement, len(targetUsers)*2+1)
	dearBlockElements[0] = slack.NewTextBlockObject("plain_text", "to:", false, false)
	i := 1
	for _, user := range targetUsers {
		dearBlockElements[i] = slack.NewImageBlockElement(user.Profile.Image48, user.RealName)
		i++
		dearBlockElements[i] = slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<@%s>", user.ID), false, false)
		i++
	}
	dearBlock := slack.NewContextBlock(
		"context",
		dearBlockElements...,
	)

	messageBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject("plain_text", message+" +"+point, false, false),
		nil,
		nil,
	)
	return []slack.Block{dearBlock, messageBlock}
}

func PostSlackMessageUsecase(slackToken string, senderSlackUserId string, slackUserIDs []string, message string, point string) error {
	// 送信者のユーザー情報を取得する
	senderUser, err := gateway.SlackGetUserInfo(slackToken, senderSlackUserId)
	if err != nil {
		util.Logger.Error(err.Error())
		return err
	}

	// 送信先ユーザーの情報を取得する
	var targetUsers []*slack.User
	for _, slackUserID := range slackUserIDs {
		user, err := gateway.SlackGetUserInfo(slackToken, slackUserID)
		if err != nil {
			util.Logger.Error(err.Error())
			return err
		}
		targetUsers = append(targetUsers, user)
	}

	blocks := generateMessageBlocks(senderUser, targetUsers, message, point)

	err = gateway.SlackPostMessage(slackToken, blocks, senderUser)
	if err != nil {
		util.Logger.Error(err.Error())
		return err
	}
	return nil
}
