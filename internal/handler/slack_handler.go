package handler

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
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
		c.AbortWithStatusJSON(401, gin.H{"message": err})
		return
	}

	s, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		util.Logger.Error(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"message": err})
		return
	}

	switch s.Command {
	case "/unipos":
		err := usecase.OpenSlackModalUsecase(c, h.Db, h.Token, s.TriggerID, s.UserID)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"message": err})
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
		c.AbortWithStatusJSON(401, gin.H{"message": err})
		return
	}

	// パラメーターの値を取得する
	var i slack.InteractionCallback
	err = json.Unmarshal([]byte(c.Request.FormValue("payload")), &i)
	if err != nil {
		util.Logger.Error(err.Error())
		c.AbortWithStatusJSON(401, gin.H{"message": err})
		return
	}
	senderSlackUserId := i.User.ID
	slackUserIDs := i.View.State.Values["Members"]["member"].SelectedUsers
	pointStr := i.View.State.Values["Point"]["point"].Value
	message := i.View.State.Values["Message"]["message"].Value

	point, err := strconv.Atoi(pointStr)
	if err != nil {
		util.Logger.Error(err.Error())
		c.AbortWithStatusJSON(400, gin.H{"message": err})
		return
	}

	// カードを作成する
	err = usecase.CreateCardUsecase(c, h.Db, senderSlackUserId, slackUserIDs, point, message)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err})
		return
	}

	// slackメッセージを送信する
	err = usecase.PostSlackMessageUsecase(h.Token, senderSlackUserId, slackUserIDs, message, pointStr)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err})
		return
	}
}
