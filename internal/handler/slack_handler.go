package handler

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/youichiro/go-slack-my-unipos/internal/util"
)

type SlackHandler struct {
	SigninSecret string
	Token        string
}

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

func (h SlackHandler) HandleSlash(c *gin.Context) {
	err := util.VerifySigningSecret(c, h.SigninSecret)
	if err != nil {
		c.IndentedJSON(401, gin.H{"message": err})
	}

	s, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.IndentedJSON(500, gin.H{"message": err})
	}

	switch s.Command {
	case "/unipos":
		api := slack.New(h.Token)
		modalRequest := generateModalRequest()
		_, err = api.OpenView(s.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
		}
	default:
		c.IndentedJSON(500, gin.H{"message": "hg"})
	}
}

func (h SlackHandler) HandleModal(c *gin.Context) {
	err := util.VerifySigningSecret(c, h.SigninSecret)
	if err != nil {
		c.IndentedJSON(401, gin.H{"message": err})
	}

	var i slack.InteractionCallback
	err = json.Unmarshal([]byte(c.Request.FormValue("payload")), &i)
	if err != nil {
		c.IndentedJSON(401, gin.H{"message": err})
	}

	// members := i.View.State.Values["Members"]["member"].SelectedUsers
	// member := members[0]
	point := i.View.State.Values["Point"]["point"].Value
	message := i.View.State.Values["Message"]["message"].Value

	msg := fmt.Sprintf("point: %s, message: %s", point, message)

	api := slack.New(h.Token)
	_, _, err = api.PostMessage(
		os.Getenv("SLACK_UNIPOS_CHANNEL_ID"),
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAttachments(),
	)
	if err != nil {
		c.IndentedJSON(401, gin.H{"message": err})
	}
}
