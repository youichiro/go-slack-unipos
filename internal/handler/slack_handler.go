package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/youichiro/go-slack-my-unipos/internal/gateway"
	"github.com/youichiro/go-slack-my-unipos/internal/usecase"
	"github.com/youichiro/go-slack-my-unipos/internal/util"
)

type SlackHandler struct {
	Db           *sql.DB
	SigninSecret string
	Token        string
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
		err := gateway.SlackOpenModal(h.Token, s.TriggerID)
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
	senderSlackUserId := i.User.ID
	slackUserIDs := i.View.State.Values["Members"]["member"].SelectedUsers
	pointStr := i.View.State.Values["Point"]["point"].Value
	message := i.View.State.Values["Message"]["message"].Value

	point, err := strconv.Atoi(pointStr)
	if err != nil {
		log.Error().Msg(err.Error())
		c.IndentedJSON(400, gin.H{"message": err})
		return
	}

	// カードを作成する
	err = usecase.CreateCardUsecase(c, h.Db, senderSlackUserId, slackUserIDs, point, message)
	if err != nil {
		log.Error().Msg(err.Error())
		c.IndentedJSON(500, gin.H{"message": err})
		return
	}

	mentionMsg := ""
	for _, slackUserID := range slackUserIDs {
		mentionMsg += "<@" + slackUserID + ">"
	}
	msg := fmt.Sprintf("from: <@%s>, to: %s, point: %d, message: %s", senderSlackUserId, mentionMsg, point, message)

	gateway.SlackPostMessage(h.Token, msg)
	if err != nil {
		log.Error().Msg(err.Error())
		c.IndentedJSON(401, gin.H{"message": err})
		return
	}
}
