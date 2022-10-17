package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

func (h SlackHandler) HandleSlash(c *gin.Context) {
	err := util.VerifySlackSigningSecret(c, h.SigninSecret)
	if err != nil {
		log.Warn().Msg(err.Error())
		c.IndentedJSON(401, gin.H{"message": err})
		return
	}

	s, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		log.Error().Msg(err.Error())
		c.IndentedJSON(500, gin.H{"message": err})
		return
	}

	switch s.Command {
	case "/unipos":
		api := slack.New(h.Token)
		modalRequest := generateModalRequest()
		_, err = api.OpenView(s.TriggerID, modalRequest)
		if err != nil {
			log.Error().Msg(err.Error())
			c.IndentedJSON(500, gin.H{"message": err})
			return
		}
	default:
		c.IndentedJSON(204, gin.H{"message": "Command not registered: " + s.Command})
		return
	}
}

func (h SlackHandler) HandleModal(c *gin.Context) {
	err := util.VerifySlackSigningSecret(c, h.SigninSecret)
	if err != nil {
		log.Warn().Msg(err.Error())
		c.IndentedJSON(401, gin.H{"message": err})
		return
	}

	var i slack.InteractionCallback
	err = json.Unmarshal([]byte(c.Request.FormValue("payload")), &i)
	if err != nil {
		log.Error().Msg(err.Error())
		c.IndentedJSON(401, gin.H{"message": err})
		return
	}

	sendUserId := i.User.ID
	userIDs := i.View.State.Values["Members"]["member"].SelectedUsers
	userIDsMsg := ""
	for _, userID := range userIDs {
		userIDsMsg += "<@" + userID + ">"
	}
	pointStr := i.View.State.Values["Point"]["point"].Value
	message := i.View.State.Values["Message"]["message"].Value

	point, err := strconv.Atoi(pointStr)
	if err != nil {
		log.Error().Msg(err.Error())
		c.IndentedJSON(400, gin.H{"message": err})
		return
	}

	// TODO: ここでポイントを消化するUsecaseを呼び出す

	msg := fmt.Sprintf("from: <@%s>, to: %s, point: %d, message: %s", sendUserId, userIDsMsg, point, message)

	api := slack.New(h.Token)
	_, _, err = api.PostMessage(
		os.Getenv("SLACK_UNIPOS_CHANNEL_ID"),
		slack.MsgOptionText(msg, false),
		slack.MsgOptionEnableLinkUnfurl(),
	)
	if err != nil {
		log.Error().Msg(err.Error())
		c.IndentedJSON(401, gin.H{"message": err})
		return
	}
}
