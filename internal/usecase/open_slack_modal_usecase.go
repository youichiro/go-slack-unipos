package usecase

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/youichiro/go-slack-my-unipos/internal/gateway"
	"github.com/youichiro/go-slack-my-unipos/internal/models"
	"github.com/youichiro/go-slack-my-unipos/internal/query"
	"github.com/youichiro/go-slack-my-unipos/internal/util"
)

func findMember(ctx *gin.Context, db *sql.DB, slackUserID string) (*models.Member, error) {
	boil.DebugMode = true

	isExists, err := models.Members(qm.Where("slack_user_id = ?", slackUserID)).Exists(ctx, db)
	if err != nil {
		util.Logger.Error(err.Error())
		return nil, err
	}
	if isExists {
		member, err := models.Members(qm.Where("slack_user_id = ?", slackUserID)).One(ctx, db)
		if err != nil {
			util.Logger.Error(err.Error())
			return nil, err
		}
		return member, nil
	} else {
		return nil, nil
	}
}

func generateModalViewRequest(remainPoint int) slack.ModalViewRequest {
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
	pointHint := slack.NewTextBlockObject("plain_text", "1人あたりのポイントを選んでください.\nあなたの今週おくれるポイントは "+strconv.Itoa(remainPoint)+" ptです.", false, false)
	point := slack.NewInputBlock("Point", pointLabel, pointHint, pointElement)

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

func OpenSlackModalUsecase(ctx *gin.Context, db *sql.DB, token string, triggerID string, slackUserID string) error {
	member, err := findMember(ctx, db, slackUserID)
	if err != nil {
		return err
	}

	remainPoint, _ := strconv.Atoi(os.Getenv("MAX_WEEK_POINT"))
	if member != nil {
		remainPoint, err = query.RemainPointQuery(ctx, db, member.ID)
		if err != nil {
			return err
		}
	}

	modalViewRequest := generateModalViewRequest(remainPoint)
	err = gateway.SlackOpenModal(token, triggerID, modalViewRequest)
	if err != nil {
		util.Logger.Error(err.Error())
		return err
	}
	return nil
}
