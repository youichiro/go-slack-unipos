package gateway

import (
	"github.com/slack-go/slack"
	"github.com/youichiro/go-slack-my-unipos/internal/util"
)

func generateModalRequest() slack.ModalViewRequest {
	titleText := slack.NewTextBlockObject("plain_text", "Unipos", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "とじる", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "おくる", false, false)

	memberSelectLabel := slack.NewTextBlockObject("plain_text", "誰に送りますか？", false, false)
	memberSelectPlaceholder := slack.NewTextBlockObject("plain_text", "メンバーを選択してください", false, false)
	memberSelectOptions := slack.NewOptionsSelectBlockElement(slack.MultiOptTypeUser, memberSelectPlaceholder, "member")
	memberSelect := slack.NewInputBlock("Members", memberSelectLabel, nil, memberSelectOptions)

	pointLabel := slack.NewTextBlockObject("plain_text", "ポイント", false, false)
	pointPlaceholder := slack.NewTextBlockObject("plain_text", "39", false, false)
	pointElement := slack.NewPlainTextInputBlockElement(pointPlaceholder, "point")
	pointElement.MaxLength = 3
	point := slack.NewInputBlock("Point", pointLabel, nil, pointElement)

	messageLabel := slack.NewTextBlockObject("plain_text", "メッセージ本文", false, false)
	messagePlaceholder := slack.NewTextBlockObject("plain_text", "感謝の気持ちを言葉にしよう！", false, false)
	messageElement := slack.NewPlainTextInputBlockElement(messagePlaceholder, "message")
	messageElement.Multiline = true
	message := slack.NewInputBlock("Message", messageLabel, nil, messageElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			memberSelect,
			point,
			message,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func SlackOpenModal(token string, triggerID string) error {
	api := slack.New(token)
	modalRequest := generateModalRequest()
	_, err := api.OpenView(triggerID, modalRequest)
	if err != nil {
		util.Logger.Error(err.Error())
		return err
	}
	return nil
}
