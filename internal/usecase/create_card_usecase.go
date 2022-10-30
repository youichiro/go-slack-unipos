package usecase

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/youichiro/go-slack-my-unipos/internal/models"
	"github.com/youichiro/go-slack-my-unipos/internal/util"
)

func findOrCreareMember(ctx *gin.Context, db *sql.DB, slackUserId string) (*models.Member, error) {
	// メンバーを取得する. もし存在しない場合は作成する.
	isExists, err := models.Members(qm.Where("slack_user_id = ?", slackUserId)).Exists(ctx, db)
	if err != nil {
		return nil, err
	}

	var member *models.Member
	if isExists {
		member, err = models.Members(qm.Where("slack_user_id = ?", slackUserId)).One(ctx, db)
		if err != nil {
			util.Logger.Error(err.Error())
			return nil, err
		}
	} else {
		member = &models.Member{SlackUserID: slackUserId}
		err = member.Insert(ctx, db, boil.Infer())
		if err != nil {
			util.Logger.Error(err.Error())
			return nil, err
		}
		util.Logger.Debug("created a new member: " + member.SlackUserID)
	}
	return member, nil
}

func CreateCardUsecase(ctx *gin.Context, db *sql.DB, senderSlackUserId string, distinationSlackUserIds []string, point int, message string) error {
	boil.DebugMode = true

	// 送信元のメンバーを取得する
	senderMember, err := findOrCreareMember(ctx, db, senderSlackUserId)
	if err != nil {
		util.Logger.Error(err.Error())
		return err
	}

	// カードを取得する
	cards, err := models.Cards(qm.Where("sender_member_id = ?", senderMember.ID)).All(ctx, db)
	if err != nil {
		util.Logger.Error(err.Error())
		return err
	}

	// 残pointを取得する
	remainPoint := 400
	if len(cards) > 0 {
		for _, card := range cards {
			remainPoint = remainPoint - card.Point
		}
	}
	for i := 0; i < len(distinationSlackUserIds); i++ {
		remainPoint = remainPoint - point
	}
	util.Logger.Debug(fmt.Sprintf("member_id: %d, remainPoint: %d", senderMember.ID, remainPoint))

	// もしポイントが足りなかったらエラーにする
	if remainPoint < 0 {
		return fmt.Errorf("not enough points")
	}

	// 送信先のメンバーの数だけカードを作成する
	for _, distinationSlackUserId := range distinationSlackUserIds {
		distinationMember, err := findOrCreareMember(ctx, db, distinationSlackUserId)
		if err != nil {
			return err
		}
		newCard := models.Card{
			SenderMemberID:      senderMember.ID,
			DistinationMemberID: distinationMember.ID,
			Point:               point,
			Message:             message,
		}
		err = newCard.Insert(ctx, db, boil.Infer())
		if err != nil {
			util.Logger.Error(err.Error())
			return err
		}
		util.Logger.Debug("A new card is cerated!")
	}

	return nil
}
